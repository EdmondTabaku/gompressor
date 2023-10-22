package compressor

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const dictDelim = `\d\`
const occDelim = `\o%d\`
const extDelim = `\e%s\`

type Compressor interface {
	Compress(content string) string
	Decompress(content string) string
}

type Base struct {
	ext string
}

type kv struct {
	Key   string
	Value int
}

func NewCompressorBase(ext string) *Base {
	return &Base{ext: ext}
}

// Compress compresses the file content
func (c *Base) Compress(file *os.File) (string, error) {
	occurrencesMap := make(map[string]int)
	var occurrences = make(chan kv)
	const maxGoroutines = 20
	sem := make(chan struct{}, maxGoroutines) // semaphore pattern
	var wg sync.WaitGroup
	linesRead := 0
	compressed := ""
	contentToProcess := ""

	// start a goroutine to save the occurrences in a map
	go func() {
		for occ := range occurrences {
			occurrencesMap[occ.Key]++
		}
	}()

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linesRead++
		contentToProcess += scanner.Text()
		if compressed != "" {
			compressed += "\n" + scanner.Text()
		} else {
			compressed += scanner.Text()
		}

		if linesRead >= 10 {
			sem <- struct{}{} // will block when the channel is full
			wg.Add(1)
			go findOccurrences(contentToProcess, &wg, sem, occurrences) // start a goroutine to find the occurrences
			contentToProcess = ""
			linesRead = 0
		}
	}

	// process the last lines
	if linesRead > 0 {
		sem <- struct{}{} // will block when the channel is full
		wg.Add(1)
		go findOccurrences(contentToProcess, &wg, sem, occurrences)
	}

	wg.Wait()
	close(occurrences)

	sizeDict := make(map[string]int)

	// save the occurrences that are worth to be replaced by a variable
	for word, occ := range occurrencesMap {
		if occ > 2 &&
			calculateSize(word, occ) > calculateSize(fmt.Sprintf(occDelim, len(sizeDict)), occ)+calculateSize(word, 1) {
			sizeDict[word] = calculateSize(word, occ) - calculateSize(fmt.Sprintf(occDelim, len(sizeDict)), occ) - calculateSize(word, 1)
		}
	}

	// sort the occurrences by size
	sortedDict := sortDict(sizeDict)

	variablesAdded := 0
	variablesHolder := dictDelim

	// replace the occurrences by variables
	for _, oc := range sortedDict {
		if strings.Contains(compressed, oc.Key) {
			varName := fmt.Sprintf(occDelim, variablesAdded)
			compressed = strings.Replace(compressed, oc.Key, varName, -1)
			variablesHolder = dictDelim + oc.Key + variablesHolder
			variablesAdded++
		}
	}

	// add the variables at the beginning of the content
	compressed = variablesHolder + compressed

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return fmt.Sprintf(extDelim, c.ext) + compressed, nil
}

// Decompress decompresses the file content
func (c *Base) Decompress(content string) (string, string, error) {
	strings.Count(content, dictDelim)
	decompressed := content
	variable := ""

	extRe, err := regexp.Compile(`\\e(.*?)\\`)
	if err != nil {
		return "", "", err
	}

	// Extract the extension from the content
	ext := extRe.FindStringSubmatch(decompressed)
	if ext != nil && len(ext) > 1 {
		decompressed = strings.Replace(decompressed, ext[0], "", 1)
	} else {
		return "", "", fmt.Errorf("Gozip file is corrupted")
	}

	re, err := regexp.Compile(fmt.Sprintf(`\%s\(.*?)\%s\`, dictDelim, dictDelim))
	if err != nil {
		return "", "", err
	}

	// Replace all variables by their content
	for i := strings.Count(content, dictDelim) - 2; i >= 0; i-- {
		variable = fmt.Sprintf(occDelim, i)
		match := re.FindStringSubmatch(decompressed)

		if match != nil && len(match) > 1 {
			// Extract the variable content
			variableContent := match[1]

			// Remove the first delimiter and the variable content
			decompressed = strings.Replace(decompressed, `\d\`+variableContent, "", 1)

			// Replace the variable with the variable content
			decompressed = strings.Replace(decompressed, variable, variableContent, -1)
		} else {
			fmt.Println("No match found")
		}

	}

	return strings.Replace(decompressed, dictDelim, "", 1), ext[1], nil
}

// findOccurrences finds the occurrences in the content and sends them to the occurrences channel
func findOccurrences(content string, wg *sync.WaitGroup, sem chan struct{}, occurrences chan kv) {
	defer wg.Done()
	defer func() { <-sem }() // release the semaphore
	var length int

	if len(content) > 200 {
		length = 100
	} else {
		length = len(content) / 2
	}

	if length > 3 { // minimum length of an occurrence
		for offset := 3; offset < length; offset++ { // offset is the length of the occurrence
			for i := 0; i < len(content)-offset; i = i + offset { // i is the position of the occurrence
				if i >= 3 && len(content) >= i+offset+3 && !containsDelimiters(content[i-3:i+offset+3]) { // check if the occurrence is not surrounded by delimiters
					// send the occurrence to the channel
					occurrences <- kv{Key: content[i : i+offset]}
				}
			}
		}
	}
}

// calculateSize calculates the size of the occurrence
func calculateSize(word string, occurrences int) int {
	return len([]byte(word)) * occurrences
}

// containsDelimiters checks if the word contains delimiters
func containsDelimiters(word string) bool {
	return strings.Contains(word, dictDelim) || strings.Contains(word, `\o`)
}

// sortDict sorts the dictionary by value
func sortDict(dict map[string]int) []kv {
	var pairList []kv
	for k, v := range dict {
		pairList = append(pairList, kv{k, v})
	}

	sort.Slice(pairList, func(i, j int) bool {
		return pairList[i].Value > pairList[j].Value
	})

	return pairList
}
