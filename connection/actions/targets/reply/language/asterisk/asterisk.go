package asterisk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	rangeSeparator = "-"
	exprSeparator = ","
)

func ProcessAsterisk(fullExpression string, buff *[]byte) (*[]byte, error) {
	var finishedBuff []byte
	expressions := strings.Split(fullExpression, exprSeparator)
	for _, expr := range expressions {
		processedExpr, err := processExpression(expr, buff)
		if err != nil {
			return &finishedBuff, err
		}
		finishedBuff = append(finishedBuff, processedExpr ...)
	}
	return &finishedBuff, nil
}

func processExpression(expression string, buff *[]byte) ([]byte, error) {
	var cutBuff []byte
	if strings.Contains(expression, rangeSeparator) {
		return processRange(expression, buff)
	} else {
		single, err := processSingle(expression, buff)
		if err != nil {
			return cutBuff, err
		}
		return append(cutBuff, single), nil
	}
}

func processSingle(expression string, buff *[]byte) (byte, error) {
	index, err := strconv.Atoi(expression)
	if err != nil {
		var single byte
		return single, err
	}
	return (*buff)[index], nil
}

func processRange(rangeExpression string, buff *[]byte) ([]byte, error) {
	rangeValues := strings.Split(rangeExpression, rangeSeparator)
	if len(rangeValues) != 2 {
		return *buff, errors.New(fmt.Sprintf("invalid range in expression-%s", rangeExpression))
	}
	builtRange, err := validateRange(rangeValues, buff)
	if err != nil {
		return *buff, err
	}
	return (*buff)[builtRange[0]:builtRange[1]+1], nil
}

func validateRange(splitRangeExpression []string, buff *[]byte) ([]int, error) {
	var rangeNums []int
	isLast := false
	for _, val := range splitRangeExpression {
		if len(val) == 0 {
			if isLast {
				rangeNums = append(rangeNums, len(*buff))
				isLast = true
				continue
			} else {
				rangeNums = append(rangeNums, 0)
				isLast = true
				continue
			}
		}
		index, err := strconv.Atoi(val)
		if err != nil {
			return rangeNums, err
		}
		if index < 0 || index >= len(*buff) {
			return rangeNums, errors.New("index outside of message length")
		}
		rangeNums = append(rangeNums, index)
		isLast = true
	}
	return rangeNums, nil
}