package ftp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type FtpState struct {
	conn          *net.UDPConn
	dst           *net.UDPAddr
	receiveBuffer map[uint32]string
}

var commands []Command = []Command{
	{CommandName: CommandName{Name: "connect"}, Desc: "connect to remote myftp"},
	{CommandName: CommandName{Name: "put"}, Desc: "send file"},
	{CommandName: CommandName{Name: "get"}, Desc: "connect to remote myftp"},
	{CommandName: CommandName{Name: "quit"}, Desc: "exit myftp"},
	{CommandName: CommandName{Name: "?"}, Desc: "print help information"},
}

func Run(dstAddr string, port int) {
	state := &FtpState{}
	state.receiveBuffer = make(map[uint32]string)

	for {
		cmd, err := state.waitForCommand()

		if err != nil {
			log.Println(err)
			continue
		}

		if v, ok := cmd["name"]; ok {
			switch v {
			case "connect":
				err = state.connect(dstAddr)
			case "put":
				err = state.put(cmd["filepath"])
			case "get":
				err = state.get(cmd["remotepath"])
			case "quit":
				log.Println("Quitting program.")
				os.Exit(0)
			case "help":
				state.displayHelp()
			}
			if err != nil {
				log.Println(err)
				continue
			}
		} else {
			log.Fatalln("Unknown error: command name not found")
		}
	}
}

func (s *FtpState) waitForCommand() (map[string]string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("myftp> ")
	line, err := reader.ReadString('\n')

	if err != nil {
		log.Fatalln(err)
	}

	lineSegments := strings.Split(strings.Trim(line, "\n"), " ")

	if len(lineSegments) == 0 {
		lineSegments = []string{"?"}
	}

	cmd := lineSegments[0]

	if cmd == "connect" || cmd == "c" {
		return map[string]string{"name": cmd}, nil
	} else if cmd == "put" || cmd == "p" {
		if len(lineSegments) != 2 {
			return nil, errors.New("missing file path")
		}
		return map[string]string{"name": cmd, "filepath": lineSegments[1]}, nil
	} else if cmd == "get" || cmd == "g" {
		if len(lineSegments) != 2 {
			return nil, errors.New("missing remote path")
		}
		return map[string]string{"name": cmd, "remotepath": lineSegments[1]}, nil
	} else if cmd == "quit" || cmd == "q" {
		return map[string]string{"quit": cmd}, nil
	} else {
		return map[string]string{"name": "help"}, nil
	}
}

func (s *FtpState) connect(dstAddr string) error {
	// Resolve the UDP address
	dst, err := net.ResolveUDPAddr("udp", dstAddr)
	if err != nil {
		return err
	}

	// Dial to the server
	conn, err := net.DialUDP("udp", nil, dst)
	if err != nil {
		return err
	}

	s.conn = conn
	s.dst = dst
	return nil
}

func (s *FtpState) get(remotePath string) error {
	if s.conn == nil || s.dst == nil {
		return errors.New("connection not made yet")
	}

	// Send the request to the server
	_, err := s.conn.Write([]byte(remotePath))
	if err != nil {
		return err
	}

	buffer := make([]byte, 512) // Adjust buffer size as needed
	file, err := os.Create(filepath.Base(remotePath))
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		n, _, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			if err == io.EOF {
				break // End of file
			}
			return err
		}
		_, err = file.Write(buffer[:n])
		if err != nil {
			return err
		}
	}

	fmt.Println("File received successfully")
	return nil
}

func (s *FtpState) put(filePath string) error {
	if s.conn == nil || s.dst == nil {
		return errors.New("connection not made yet")
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file contents into a buffer
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Send the file data to the server
	_, err = s.conn.Write(buffer)
	if err != nil {
		return err
	}

	fmt.Println("File sent successfully")
	return nil
}

func (s *FtpState) displayHelp() {
	fmt.Println("Commands may be abbreviated. Commands are:")
	for _, command := range commands {
		fmt.Printf("(%s)%-10s %s\n",
			string(command.Name[0]), command.Name[1:], command.Desc)
	}
}
