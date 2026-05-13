package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/pushtisonawala/jaeger-skills-poc/skills"
)

func main() {
	dir := flag.String("skills-dir", "./examples", "Directory to load skills from")
	dryRun := flag.Bool("dry-run", false, "Validate skills without registering")
	flag.Parse()

	fmt.Println("+------------------------------------------+")
	fmt.Println("|     Jaeger Skills Engine - PoC           |")
	fmt.Println("|     BYOA Skills Framework Prototype      |")
	fmt.Println("+------------------------------------------+")
	fmt.Println()

	loaded, loadErrs := skills.LoadFromDir(*dir)
	if len(loadErrs) > 0 {
		for _, e := range loadErrs {
			color.Red("  LOAD ERROR: %v", e)
		}
	}
	fmt.Printf("Discovered %d skill file(s) in %s\n\n", len(loaded), *dir)

	registry := skills.NewRegistry()
	hasErrors := len(loadErrs) > 0

	for _, skill := range loaded {
		result := skills.Validate(skill)

		if result.Valid {
			color.Green("  OK %-35s [VALID]", skill.Name)
		} else {
			color.Red("  XX %-35s [INVALID]", skill.Name)
			hasErrors = true
		}

		for _, e := range result.Errors {
			color.Red("      ERROR:   %s", e)
		}
		for _, w := range result.Warnings {
			color.Yellow("      WARNING: %s", w)
		}

		if result.Valid && !*dryRun {
			if err := registry.Register(skill); err != nil {
				color.Red("      CONFLICT: %v", err)
				hasErrors = true
			}
		}
	}

	fmt.Println()

	if *dryRun {
		fmt.Println("-- Dry-run mode: validation only, nothing registered --")
	} else {
		fmt.Println("-- Registry --")
		fmt.Print(registry.Summary())
	}

	fmt.Println()
	if hasErrors {
		color.Red("Result: FAILED - fix errors above before deploying skills")
		os.Exit(1)
	}

	color.Green("Result: OK - all skills loaded successfully")
}
