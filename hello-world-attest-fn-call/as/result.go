package as

type Result struct {
	Success bool
	Error   string
	Output  any
}

func WriteOutput(output any) uint64 {
	return WriteResult(Result{
		Success: true,
		Output:  output,
	})
}

func WriteError(err string) uint64 {
	return WriteResult(Result{
		Success: false,
		Error:   err,
	})
}

func WriteResult(result Result) uint64 {
	outputData, err := Marshal(result)
	if err != nil {
		panic("Fatal error: could not marshal output data")
	}
	return ShareWithHost(outputData)
}
