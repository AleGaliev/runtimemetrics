package filestore

import (
	"bufio"
	"fmt"
	"os"
)

type Producer struct {
	file *os.File // файл для записи
}

type Consumer struct {
	file *os.File
	// заменяем Reader на Scanner
	scanner *bufio.Scanner
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file: file,
	}, nil
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (p *Producer) Close() error {
	// закрываем файл
	return p.file.Close()
}

func WriteMetrics(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("could not open metrics file: %w", err)
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Errorf("could not write metrics to file: %w", err)
	}
	file.Close()

	return nil
}

func ReadMetrics(filename string) ([]byte, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open metrics file: %w", err)
	}

	info, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read metrics file: %w", err)
	}
	if info.Size() == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("could not read metrics from file: %w", scanner.Err())
	}

	data := scanner.Bytes()

	file.Close()

	return data, nil
}

func (c *Consumer) Close() error {
	// закрываем файл
	return c.file.Close()
}
