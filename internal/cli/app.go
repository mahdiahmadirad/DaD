package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	framework "github.com/mahdiahmadirad/DaD"
	core "github.com/mahdiahmadirad/DaD/internal/dad"
)

const Version = "0.1.0"

type globalOptions struct {
	root    string
	format  string
	quiet   bool
	noColor bool
	dryRun  bool
	help    bool
	version bool
}

type outcome struct {
	command     string
	code        int
	data        any
	diagnostics []core.Diagnostic
	text        string
}

type envelope struct {
	SchemaVersion string            `json:"schema_version"`
	Command       string            `json:"command"`
	Success       bool              `json:"success"`
	Data          any               `json:"data"`
	Diagnostics   []core.Diagnostic `json:"diagnostics"`
}

func Run(args []string, stdout, stderr io.Writer, environment []string) int {
	options, remaining, parseOutcome := parseGlobals(args)
	if parseOutcome != nil {
		emit(*parseOutcome, options, stdout, stderr)
		return parseOutcome.code
	}
	if options.version {
		result := outcome{
			command: "version", code: core.ExitOK,
			data: map[string]string{"version": Version},
			text: "dad " + Version + "\n",
		}
		emit(result, options, stdout, stderr)
		return result.code
	}
	if len(remaining) == 0 {
		result := outcome{
			command: "help", code: core.ExitOK,
			data: map[string]string{"help": generalHelp()},
			text: generalHelp(),
		}
		if !options.help {
			result.code = core.ExitUsage
			result.diagnostics = []core.Diagnostic{{
				Code: "DAD-USAGE-001", Severity: "error",
				Message: "A command is required.",
			}}
		}
		emit(result, options, stdout, stderr)
		return result.code
	}

	command := strings.ToLower(remaining[0])
	commandArgs := remaining[1:]
	if command == "help" {
		options.help = true
		if len(commandArgs) > 0 {
			command, commandArgs = strings.ToLower(commandArgs[0]), commandArgs[1:]
		} else {
			result := outcome{
				command: "help", code: core.ExitOK,
				data: map[string]string{"help": generalHelp()},
				text: generalHelp(),
			}
			emit(result, options, stdout, stderr)
			return result.code
		}
	}
	if options.help {
		help, ok := commandHelp(command)
		if !ok {
			result := usage(command, fmt.Sprintf("Unknown command %q.", command))
			emit(result, options, stdout, stderr)
			return result.code
		}
		result := outcome{
			command: command, code: core.ExitOK,
			data: map[string]string{"help": help}, text: help,
		}
		emit(result, options, stdout, stderr)
		return result.code
	}
	if options.dryRun && command != "init" && command != "new" {
		result := usage(command, "--dry-run is valid only for init and new.")
		emit(result, options, stdout, stderr)
		return result.code
	}

	cwd, err := os.Getwd()
	if err != nil {
		result := failure(command, core.ExitIO, "DAD-IO-100", err.Error())
		emit(result, options, stdout, stderr)
		return result.code
	}
	env := environmentMap(environment)
	result := dispatch(command, commandArgs, options, cwd, env)
	emit(result, options, stdout, stderr)
	return result.code
}

func dispatch(
	command string,
	args []string,
	options globalOptions,
	cwd string,
	environment map[string]string,
) outcome {
	switch command {
	case "init":
		return runInit(args, options, cwd)
	case "new":
		return withRoot(command, options, cwd, environment, func(root string) outcome {
			return runNew(root, args, options)
		})
	case "list":
		return withRoot(command, options, cwd, environment, func(root string) outcome {
			return runList(root, args)
		})
	case "status":
		return withRoot(command, options, cwd, environment, func(root string) outcome {
			return runStatus(root, args)
		})
	case "context":
		return withRoot(command, options, cwd, environment, func(root string) outcome {
			return runContext(root, args)
		})
	case "check":
		return withRoot(command, options, cwd, environment, func(root string) outcome {
			return runCheck(root, args)
		})
	case "prompt":
		return runPrompt(args)
	default:
		return usage(command, fmt.Sprintf("Unknown command %q.", command))
	}
}

func withRoot(
	command string,
	options globalOptions,
	cwd string,
	environment map[string]string,
	run func(string) outcome,
) outcome {
	root, err := core.ResolveRoot(options.root, environment, cwd)
	if err != nil {
		return failure(command, core.ExitRoot, "DAD-ROOT-001", err.Error())
	}
	return run(root)
}

func runInit(args []string, options globalOptions, cwd string) outcome {
	positionals, _, _, err := parseCommandArgs(args, nil, nil)
	if err != nil || len(positionals) > 1 {
		if err == nil {
			err = fmt.Errorf("init accepts at most one path")
		}
		return usage("init", err.Error())
	}
	path := ""
	if len(positionals) == 1 && options.root != "" {
		return usage("init", "Use either an init path or --root, not both.")
	}
	if len(positionals) == 1 {
		path = positionals[0]
	} else if options.root != "" {
		path = options.root
	}
	root, err := core.ResolveInitTarget(path, cwd)
	if err != nil {
		return failure("init", core.ExitRoot, "DAD-ROOT-001", err.Error())
	}
	files, err := framework.SupportFiles()
	if err != nil {
		return failure("init", core.ExitInternal, "DAD-INTERNAL-001", err.Error())
	}
	result, diagnostics, code := core.Init(root, files, options.dryRun)
	var text strings.Builder
	for _, path := range result.Created {
		if options.dryRun {
			fmt.Fprintf(&text, "would create %s\n", path)
		} else {
			fmt.Fprintf(&text, "created %s\n", path)
		}
	}
	for _, path := range result.Unchanged {
		fmt.Fprintf(&text, "unchanged %s\n", path)
	}
	return outcome{
		command: "init", code: code, data: result,
		diagnostics: diagnostics, text: text.String(),
	}
}

func runNew(root string, args []string, options globalOptions) outcome {
	positionals, values, _, err := parseCommandArgs(
		args,
		map[string]bool{"title": true, "number": true},
		nil,
	)
	if err != nil || len(positionals) != 1 {
		if err == nil {
			err = fmt.Errorf("new requires exactly one type: adr, spec, or task")
		}
		return usage("new", err.Error())
	}
	documentType, ok := core.ParseType(positionals[0])
	if !ok {
		return usage("new", "Document type must be adr, spec, or task.")
	}
	title := values["title"]
	if title == "" {
		return usage("new", "--title is required.")
	}
	number := 0
	if value := values["number"]; value != "" {
		number, err = core.ParseExplicitNumber(value)
		if err != nil {
			return usage("new", err.Error())
		}
	}
	support, err := framework.SupportFiles()
	if err != nil {
		return failure("new", core.ExitInternal, "DAD-INTERNAL-001", err.Error())
	}
	templatePath := "docs/templates/" + string(documentType) + "-TEMPLATE.md"
	result, diagnostics, code := core.NewDocument(root, core.NewOptions{
		Type: documentType, Title: title, Number: number,
		DryRun: options.dryRun, Template: support[templatePath],
	})
	text := ""
	if code == core.ExitOK {
		action := "created"
		if options.dryRun {
			action = "would create"
		}
		text = fmt.Sprintf("%s %s %s\n", action, result.ID, result.Path)
	}
	return outcome{
		command: "new", code: code, data: result,
		diagnostics: diagnostics, text: text,
	}
}

func runList(root string, args []string) outcome {
	positionals, values, _, err := parseCommandArgs(
		args, map[string]bool{"status": true}, nil,
	)
	if err != nil || len(positionals) > 1 {
		if err == nil {
			err = fmt.Errorf("list accepts at most one document type")
		}
		return usage("list", err.Error())
	}
	var filterType core.DocumentType
	if len(positionals) == 1 {
		var ok bool
		filterType, ok = core.ParseType(positionals[0])
		if !ok {
			return usage("list", "Document type must be adr, spec, or task.")
		}
	}
	status := values["status"]
	if status != "" {
		valid := core.AnyStatusAllowed(status)
		if filterType != "" {
			valid = core.StatusAllowed(filterType, status)
		}
		if !valid {
			return usage("list", fmt.Sprintf("Invalid lifecycle status %q.", status))
		}
	}
	repository := core.LoadRepository(root)
	var documents []core.Document
	for _, document := range repository.Documents {
		if !core.StatusAllowed(document.Type, document.Status) ||
			len(repository.ByID[strings.ToUpper(document.ID)]) != 1 {
			continue
		}
		if filterType != "" && document.Type != filterType {
			continue
		}
		if status != "" && document.Status != status {
			continue
		}
		documents = append(documents, document)
	}
	code := core.ExitOK
	if hasErrors(repository.Diagnostics) {
		code = core.ExitCheck
	}
	var text strings.Builder
	for _, document := range documents {
		fmt.Fprintf(
			&text, "%s\t%s\t%s\t%s\n",
			document.ID, document.Status, document.Path, document.Title,
		)
	}
	return outcome{
		command: "list", code: code,
		data:        map[string]any{"documents": documents},
		diagnostics: repository.Diagnostics, text: text.String(),
	}
}

func runStatus(root string, args []string) outcome {
	if len(args) != 0 {
		return usage("status", "status accepts no arguments.")
	}
	result, diagnostics, code := core.Status(root)
	var text strings.Builder
	fmt.Fprintf(&text, "root %s\n", result.Root)
	for _, typeName := range []string{"ADR", "SPEC", "TASK"} {
		statuses := result.Counts[typeName]
		var names []string
		for name := range statuses {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Fprintf(&text, "%s %s %d\n", typeName, name, statuses[name])
		}
	}
	var corePaths []string
	for path := range result.CoreFiles {
		corePaths = append(corePaths, path)
	}
	sort.Strings(corePaths)
	for _, path := range corePaths {
		state := "missing"
		if result.CoreFiles[path] {
			state = "present"
		}
		fmt.Fprintf(&text, "core %s %s\n", state, path)
	}
	for _, task := range result.ActionableTasks {
		fmt.Fprintf(&text, "actionable %s %s\n", task.ID, task.Status)
	}
	for _, task := range result.BlockedTasks {
		fmt.Fprintf(&text, "blocked %s\n", task.ID)
	}
	fmt.Fprintf(
		&text, "checks errors=%d warnings=%d\n",
		result.CheckSummary.Errors, result.CheckSummary.Warnings,
	)
	return outcome{
		command: "status", code: code, data: result,
		diagnostics: diagnostics, text: text.String(),
	}
}

func runContext(root string, args []string) outcome {
	if len(args) != 1 {
		return usage("context", "context requires exactly one TASK identifier.")
	}
	result, diagnostics, code := core.Context(root, strings.ToUpper(args[0]))
	var text strings.Builder
	for _, document := range result.Documents {
		id := document.ID
		if id == "" {
			id = "-"
		}
		fmt.Fprintf(
			&text, "%s\t%s\t%s\t%s\trequired-by=%s\n",
			document.Role, id, document.Status, document.Path, document.RequiredBy,
		)
	}
	return outcome{
		command: "context", code: code, data: result,
		diagnostics: diagnostics, text: text.String(),
	}
}

func runCheck(root string, args []string) outcome {
	positionals, _, booleans, err := parseCommandArgs(
		args, nil, map[string]bool{"strict": true},
	)
	if err != nil {
		return usage("check", err.Error())
	}
	for index := range positionals {
		positionals[index] = strings.ToUpper(positionals[index])
	}
	result, diagnostics, code := core.Check(root, positionals, booleans["strict"])
	data := map[string]any{
		"summary":           result.Summary,
		"checked_documents": result.CheckedDocuments,
		"scope":             "structural-only",
	}
	text := fmt.Sprintf(
		"checked %d documents: %d errors, %d warnings\n"+
			"Structural checks do not prove architecture, behavior, freshness, "+
			"implementation completeness, or security.\n",
		len(result.CheckedDocuments), result.Summary.Errors, result.Summary.Warnings,
	)
	return outcome{
		command: "check", code: code, data: data,
		diagnostics: diagnostics, text: text,
	}
}

func runPrompt(args []string) outcome {
	if len(args) == 0 {
		return usage("prompt", "prompt requires list or show.")
	}
	switch strings.ToLower(args[0]) {
	case "list":
		if len(args) != 1 {
			return usage("prompt list", "prompt list accepts no additional arguments.")
		}
		prompts := framework.Prompts()
		var text strings.Builder
		for _, prompt := range prompts {
			fmt.Fprintf(&text, "%s\t%s\n", prompt.Name, prompt.Purpose)
		}
		return outcome{
			command: "prompt list", code: core.ExitOK,
			data: map[string]any{"prompts": prompts}, text: text.String(),
		}
	case "show":
		if len(args) != 2 {
			return usage("prompt show", "prompt show requires exactly one prompt name.")
		}
		info, content, found, err := framework.Prompt(strings.ToLower(args[1]))
		if err != nil {
			return failure("prompt show", core.ExitInternal, "DAD-INTERNAL-002", err.Error())
		}
		if !found {
			return usage("prompt show", fmt.Sprintf("Unknown prompt %q.", args[1]))
		}
		return outcome{
			command: "prompt show", code: core.ExitOK,
			data: map[string]any{
				"name": info.Name, "content": string(content), "version": Version,
			},
			text: string(content),
		}
	default:
		return usage("prompt", "prompt requires list or show.")
	}
}

func emit(result outcome, options globalOptions, stdout, stderr io.Writer) {
	core.SortDiagnostics(result.diagnostics)
	if result.data == nil {
		result.data = map[string]any{}
	}
	if result.diagnostics == nil {
		result.diagnostics = []core.Diagnostic{}
	}
	if options.format == "json" {
		encoder := json.NewEncoder(stdout)
		encoder.SetEscapeHTML(false)
		_ = encoder.Encode(envelope{
			SchemaVersion: "1", Command: result.command,
			Success: result.code == core.ExitOK, Data: result.data,
			Diagnostics: result.diagnostics,
		})
		return
	}
	if !(options.quiet && result.code == core.ExitOK && !requestedDataCommand(result.command)) {
		_, _ = io.WriteString(stdout, result.text)
	}
	for _, diagnostic := range result.diagnostics {
		_, _ = fmt.Fprintln(stderr, diagnostic.String())
	}
}

func requestedDataCommand(command string) bool {
	switch command {
	case "list", "status", "context", "check", "prompt list", "prompt show",
		"help", "version":
		return true
	default:
		return false
	}
}

func parseGlobals(args []string) (globalOptions, []string, *outcome) {
	options := globalOptions{format: "text"}
	var remaining []string
	for index := 0; index < len(args); index++ {
		argument := args[index]
		name, value, hasValue := splitOption(argument)
		switch name {
		case "--root", "--format":
			if !hasValue {
				index++
				if index >= len(args) {
					result := usage("", name+" requires a value.")
					return options, nil, &result
				}
				value = args[index]
			}
			if name == "--root" {
				options.root = value
			} else {
				options.format = value
			}
		case "--no-color":
			if hasValue {
				result := usage("", "--no-color does not accept a value.")
				return options, nil, &result
			}
			options.noColor = true
		case "--quiet":
			if hasValue {
				result := usage("", "--quiet does not accept a value.")
				return options, nil, &result
			}
			options.quiet = true
		case "--dry-run":
			if hasValue {
				result := usage("", "--dry-run does not accept a value.")
				return options, nil, &result
			}
			options.dryRun = true
		case "--help", "-h":
			if hasValue {
				result := usage("", "--help does not accept a value.")
				return options, nil, &result
			}
			options.help = true
		case "--version":
			if hasValue {
				result := usage("", "--version does not accept a value.")
				return options, nil, &result
			}
			options.version = true
		default:
			remaining = append(remaining, argument)
		}
	}
	if options.format != "text" && options.format != "json" {
		result := usage("", "--format must be text or json.")
		return options, nil, &result
	}
	return options, remaining, nil
}

func parseCommandArgs(
	args []string,
	valueOptions map[string]bool,
	boolOptions map[string]bool,
) ([]string, map[string]string, map[string]bool, error) {
	values := make(map[string]string)
	booleans := make(map[string]bool)
	var positionals []string
	for index := 0; index < len(args); index++ {
		argument := args[index]
		if !strings.HasPrefix(argument, "--") {
			positionals = append(positionals, argument)
			continue
		}
		name, value, hasValue := splitOption(argument)
		key := strings.TrimPrefix(name, "--")
		if valueOptions[key] {
			if !hasValue {
				index++
				if index >= len(args) {
					return nil, nil, nil, fmt.Errorf("%s requires a value", name)
				}
				value = args[index]
			}
			values[key] = value
			continue
		}
		if boolOptions[key] {
			if hasValue {
				return nil, nil, nil, fmt.Errorf("%s does not accept a value", name)
			}
			booleans[key] = true
			continue
		}
		return nil, nil, nil, fmt.Errorf("unknown option %s", name)
	}
	return positionals, values, booleans, nil
}

func splitOption(argument string) (name, value string, hasValue bool) {
	if index := strings.IndexByte(argument, '='); index >= 0 {
		return argument[:index], argument[index+1:], true
	}
	return argument, "", false
}

func environmentMap(values []string) map[string]string {
	result := make(map[string]string, len(values))
	for _, value := range values {
		key, item, found := strings.Cut(value, "=")
		if found {
			result[key] = item
		}
	}
	return result
}

func usage(command, message string) outcome {
	return failure(command, core.ExitUsage, "DAD-USAGE-001", message)
}

func failure(command string, code int, diagnosticCode, message string) outcome {
	return outcome{
		command: command, code: code, data: map[string]any{},
		diagnostics: []core.Diagnostic{{
			Code: diagnosticCode, Severity: "error", Message: message,
		}},
	}
}

func hasErrors(diagnostics []core.Diagnostic) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == "error" {
			return true
		}
	}
	return false
}

func generalHelp() string {
	return `DaD — Document-aware Development CLI

Usage:
  dad [global options] <command> [arguments]

Commands:
  init       Install canonical reusable DaD support files
  new        Create an ADR, SPEC, or TASK
  list       List governed documents
  status     Summarize document state and mechanical health
  context    Resolve governing context for a TASK
  check      Check document mechanics
  prompt     List or show official DaD prompts

Global options:
  --root <path>          Select the DaD repository root
  --format text|json     Select output format
  --no-color             Disable color
  --quiet                Suppress nonessential success output
  --dry-run              Preview init or new without writing
  --help                 Show help
  --version              Show version
`
}

func commandHelp(command string) (string, bool) {
	help := map[string]string{
		"init": "Usage: dad init [path] [--dry-run]\n",
		"new": "Usage: dad new <adr|spec|task> --title <title> " +
			"[--number NNNN] [--dry-run]\n",
		"list":    "Usage: dad list [adr|spec|task] [--status <status>]\n",
		"status":  "Usage: dad status\n",
		"context": "Usage: dad context <TASK-NNNN>\n",
		"check":   "Usage: dad check [document-id...] [--strict]\n",
		"prompt":  "Usage: dad prompt <list|show> [name]\n",
	}
	value, ok := help[command]
	return value, ok
}
