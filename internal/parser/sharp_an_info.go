package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	keySWGUID                 = "SW_GUID"
	keyEndianness             = "endianness"
	keyEnableEndiannessPerJob = "enable_endianness_per_job"
	keyReproducibilityDisable = "reproducibility_disable"
)

var required = []string{
	keySWGUID,
	keyEndianness,
	keyEnableEndiannessPerJob,
	keyReproducibilityDisable,
}

func ParseSharpANInfo(r io.Reader) ([]SharpANInfoRow, error) {
	result := make([]SharpANInfoRow, 0)
	params := make(map[string]string)
	scanner := bufio.NewScanner(r)
	flush := func() error {
		if len(params) == 0 {
			return nil
		}

		for _, k := range required {
			if _, ok := params[k]; !ok {
				return fmt.Errorf("missing required key %q", k)
			}
		}

		endianness, err := strconv.Atoi(params[keyEndianness])
		if err != nil {
			return fmt.Errorf("invalid endianness: %q", params[keyEndianness])
		}

		enable_endianness_per_job, err := strconv.Atoi(params[keyEnableEndiannessPerJob])
		if err != nil {
			return fmt.Errorf("invalid enable_endianness_per_job: %q", params[keyEnableEndiannessPerJob])
		}

		reproducibility_disable, err := strconv.Atoi(params[keyReproducibilityDisable])
		if err != nil {
			return fmt.Errorf("invalid reproducibility_disable: %q", params[keyReproducibilityDisable])
		}

		result = append(result, SharpANInfoRow{
			SWGUID:                 params[keySWGUID],
			Endianness:             endianness,
			EnableEndiannessPerJob: enable_endianness_per_job,
			ReproducibilityDisable: reproducibility_disable,
		})
		params = make(map[string]string)
		return nil
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case line == "":
			if err := flush(); err != nil {
				return nil, err
			}
			continue
		default:
			key, value, ok := parseINFOLine(line)
			if !ok {
				return nil, fmt.Errorf("invalid line (no key=value or empty key): %q", line)
			}
			params[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan sharp_an_info: %w", err)
	}

	// так как в конце файла может и не быть пустой строки, нужно вызвать flush() ещё разок
	if err := flush(); err != nil {
		return nil, err
	}

	return result, nil
}

func parseINFOLine(line string) (key, value string, ok bool) {
	before, after, found := strings.Cut(line, "=")
	if !found {
		return "", "", false
	}

	key = strings.TrimSpace(before)
	value = strings.TrimSpace(after)
	if key == "" {
		return "", "", false
	}

	return key, value, true
}
