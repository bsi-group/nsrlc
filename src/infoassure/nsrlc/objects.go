package main

// Holds the various objects/structs that are used in the system that don't warrant their own individual file

import (
	"github.com/voxelbrain/goptions"
)

// ##### Structs #############################################################

// Structure to store options
type Options struct {
	InputFile 	string        	`goptions:"-i, --input, obligatory, description='Input file path'"`
	OutputFile 	string        	`goptions:"-o, --output, obligatory, description='Output file path'"`
	Format 		string       	`goptions:"-f, --format, description='Format for output e.g. i (identified), u (unidentified) or a (all)'"`
	Server		string        	`goptions:"-s, --server, description='NSRL server'"`
	BatchSize	int        		`goptions:"-b, --batchsize, description='Batch size'"`
	Help    	goptions.Help	`goptions:"-h, --help, description='Show this help'"`
}

// Struct to marshal the "data" to a JSON string for the API
type JsonResult struct {
	Hash 	string		`db:"hash"`
	Exists	bool		`db:"exists"`
}