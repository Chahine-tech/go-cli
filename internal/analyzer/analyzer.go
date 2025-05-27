package parser

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"loganalyzer/internal/config"
	"loganalyzer/internal/reporter"
)

type Analyzer struct {
	config   *config.Config
	reporter *reporter.Reporter
}

func NewAnalyzer(cfg *config.Config) *Analyzer {
	return &Analyzer{
		config:   cfg,
		reporter: reporter.NewReporter(),
	}
}

func (a *Analyzer) AnalyzeAllLogs() error {
	if len(a.config.Logs) == 0 {
		return fmt.Errorf("no logs to analyze")
	}

	resultsChan := make(chan reporter.AnalysisResult, len(a.config.Logs))
	
	var wg sync.WaitGroup

	fmt.Printf("Starting analysis of %d log files...\n", len(a.config.Logs))

	for _, logConfig := range a.config.Logs {
		wg.Add(1)
		go a.analyzeLogFile(logConfig, resultsChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for result := range resultsChan {
		a.reporter.AddResult(result)
	}

	return nil
}

func (a *Analyzer) analyzeLogFile(logConfig config.LogConfig, resultsChan chan<- reporter.AnalysisResult, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Processing log: %s (%s)\n", logConfig.ID, logConfig.Path)

	if err := a.checkFileAccess(logConfig.Path); err != nil {
		result := a.handleFileError(logConfig, err)
		resultsChan <- result
		return
	}

	sleepDuration := time.Duration(50+rand.Intn(151)) * time.Millisecond
	time.Sleep(sleepDuration)

	if rand.Float32() < 0.1 {
		parseErr := NewParseError(logConfig.ID, "simulated parsing failure", fmt.Errorf("random parse error"))
		result := a.handleParseError(logConfig, parseErr)
		resultsChan <- result
		return
	}

	result := reporter.CreateSuccessResult(logConfig.ID, logConfig.Path)
	fmt.Printf("✓ Completed analysis of log: %s\n", logConfig.ID)
	resultsChan <- result
}

func (a *Analyzer) checkFileAccess(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return NewFileNotFoundError(filePath, err)
		}
		return NewFileNotFoundError(filePath, fmt.Errorf("file access error: %w", err))
	}

	if !info.Mode().IsRegular() {
		return NewFileNotFoundError(filePath, fmt.Errorf("not a regular file"))
	}

	file, err := os.Open(filePath)
	if err != nil {
		return NewFileNotFoundError(filePath, fmt.Errorf("cannot read file: %w", err))
	}
	file.Close()

	return nil
}

func (a *Analyzer) handleFileError(logConfig config.LogConfig, err error) reporter.AnalysisResult {
	var fileNotFoundErr *FileNotFoundError

	if errors.As(err, &fileNotFoundErr) {
		fmt.Printf("✗ File error for log %s: %s\n", logConfig.ID, fileNotFoundErr.Error())
		return reporter.CreateFailureResult(
			logConfig.ID,
			logConfig.Path,
			"File not found.",
			fileNotFoundErr.Error(),
		)
	}

	// Generic file error
	fmt.Printf("✗ File error for log %s: %s\n", logConfig.ID, err.Error())
	return reporter.CreateFailureResult(
		logConfig.ID,
		logConfig.Path,
		"File access error.",
		err.Error(),
	)
}

func (a *Analyzer) handleParseError(logConfig config.LogConfig, err error) reporter.AnalysisResult {
	var parseErr *ParseError

	if errors.As(err, &parseErr) {
		fmt.Printf("✗ Parse error for log %s: %s\n", logConfig.ID, parseErr.Error())
		return reporter.CreateFailureResult(
			logConfig.ID,
			logConfig.Path,
			"Parse error occurred.",
			parseErr.Error(),
		)
	}

	fmt.Printf("✗ Parse error for log %s: %s\n", logConfig.ID, err.Error())
	return reporter.CreateFailureResult(
		logConfig.ID,
		logConfig.Path,
		"Parse error occurred.",
		err.Error(),
	)
}

func (a *Analyzer) GetReporter() *reporter.Reporter {
	return a.reporter
} 