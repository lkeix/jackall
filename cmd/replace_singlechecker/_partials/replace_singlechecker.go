// It will be added to golang.org/go/analysis package

func Main(a *analysis.Analyzer) {
	log.SetFlags(0)
	log.SetPrefix(a.Name + ": ")

	analyzers := []*analysis.Analyzer{a}

	if err := analysis.Validate(analyzers); err != nil {
		log.Fatal(err)
	}

	checker.RegisterFlags()

	flag.Usage = func() {
		paras := strings.Split(a.Doc, "\n\n")
		fmt.Fprintf(os.Stderr, "%s: %s\n\n", a.Name, paras[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [-flag] [package]\n\n", a.Name)
		if len(paras) > 1 {
			fmt.Fprintln(os.Stderr, strings.Join(paras[1:], "\n\n"))
		}
		fmt.Fprintln(os.Stderr, "\nFlags:")
		flag.PrintDefaults()
	}

	analyzers = analysisflags.Parse(analyzers, false)

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if len(args) == 1 && strings.HasSuffix(args[0], ".cfg") {
		unitchecker.Run(args[0], analyzers)
		panic("unreachable")
	}

	checker.Run(args, analyzers)
	// jackall analyze the degree of each packages dependencies.
	// os.Exit(checker.Run(args, analyzers))
}

// end
