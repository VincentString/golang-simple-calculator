package main

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("String evaluation approach of a simple calculator")
	var input *string

	input = flag.String("input", "3*(3-1)*3.7/5+23*243", "input the expression to evaluate")
	flag.Parse()
	errBracket := doBracket(input)
	if errBracket != nil {
		fmt.Println(errBracket)
		return
	}
	errMultiplyDivide := doMultiplyDivide(input)
	if errMultiplyDivide != nil {
		fmt.Println(errMultiplyDivide)
		return
	}
	errPlusMinus := doPlusMinus(input)
	if errPlusMinus != nil {
		fmt.Println(errPlusMinus)
		return
	}
	fmt.Printf("Result: %v", *input)
}

func doBracket(ms *string) error {
	inBracket := regexp.MustCompile(`\(([+\-/*]?\d*\.?\d*)+\)`)
	content := inBracket.FindString(*ms)
	if content == "" {
		return nil
	}
	processedContent := new(string)
	*processedContent = content[1 : len(content)-1]
	err1 := doMultiplyDivide(processedContent)
	if err1 != nil {
		return err1
	}
	err2 := doPlusMinus(processedContent)
	if err2 != nil {
		return err2
	}
	*ms = strings.Replace(*ms, content, *processedContent, -1)
	return nil
}

func doMultiplyDivide(ms *string) error {
	multiplyDivide := regexp.MustCompile(`\d+\.?\d*[/*]\d*\.?\d*`)
	s := multiplyDivide.FindString(*ms)
	if s == "" {
		return nil
	}

	result, err := stringExpression(s)
	if err != nil {
		return err
	}
	*ms = strings.Replace(*ms, s, result, -1)
	return doMultiplyDivide(ms)
}

func doPlusMinus(ms *string) error {
	plusMinus := regexp.MustCompile(`[+\-]?\d+\.?\d*[+\-]\d*\.?\d*`)
	s := plusMinus.FindString(*ms)
	if s == "" {
		return nil
	}
	result, err := stringExpression(s)
	if err != nil {
		return err
	}
	*ms = strings.Replace(*ms, s, result, -1)
	return doPlusMinus(ms)
}

func stringExpression(s string) (string, error) {
	var result float64
	nums := []float64{}
	numStrings := make([]string, 2)
	numStrings[0] = regexp.MustCompile(`^[+\-]?\d*\.?\d*`).FindString(s)
	numStrings[1] = regexp.MustCompile(`\d*\.?\d*$`).FindString(s)

	for _, numString := range numStrings {
		num, err := strconv.ParseFloat(numString, 64)
		if err != nil {
			return "", errors.New("string can not convert to float64")
		}
		nums = append(nums, num)
	}

	ops := regexp.MustCompile(`[+\-/\*]`).FindAllString(s, -1)
	op := ops[len(ops)-1]
	switch op {
	case "+":
		result = nums[0] + nums[1]
	case "-":
		result = nums[0] - nums[1]
	case "*":
		result = nums[0] * nums[1]
	case "/":
		result = nums[0] / nums[1]
	default:
		return "", fmt.Errorf("can not recongnize the operator: %v", op)
	}
	return fmt.Sprintf("%f", result), nil
}
