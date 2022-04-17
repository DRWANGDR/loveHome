package main

import (
	"flag"
	"fmt"
	"os"
)

const Usage = `
	./v1 createchain -address ADDRESS				"create block chain"
	./v1 send -from FROM -to TO -amount AMOUNT    	"make a transaction"
	./v1 printchain									"print all blocks"
	./v1 getbalance -address ADDRESS				"get balance"
`

type CLI struct {
	//bc *BlockChain
}

func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println("Too few parameters! ", Usage)
		os.Exit(1)
	}
	createChainCmd := flag.NewFlagSet("createchain", flag.ExitOnError)
	//addBlockCmd := flag.NewFlagSet("addBlock",flag.ExitOnError)
	printCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	//addBlockCmdPara := addBlockCmd.String("data","","block info")
	createChainCmdPara := createChainCmd.String("address", "", "address data")
	getbalanceCmdPara := getbalanceCmd.String("address", "", "address data")

	fromPara := sendCmd.String("from", "", "from address data")
	toPara := sendCmd.String("to", "", "to address data")
	amountPara := sendCmd.Float64("amount", 0, "amount value")

	fmt.Println("命令：", os.Args)
	switch os.Args[1] {
	case "createchain":
		err := createChainCmd.Parse(os.Args[2:])
		CheckErr("createChainCmd", err)
		if createChainCmd.Parsed() {
			if *createChainCmdPara == "" {
				fmt.Println("Address is empty,pls check !")
				os.Exit(1)
			} else {
				cli.CreateChain(*createChainCmdPara)
			}
		}
	//case "addBlock":
	//	err := addBlockCmd.Parse(os.Args[2:])
	//	CheckErr("addBlock",err)
	//	if addBlockCmd.Parsed() {
	//		if *addBlockCmdPara == "" {
	//			fmt.Println("Data is empty,pls check !")
	//			os.Exit(1)
	//		}else {
	//			//cli.AddBlock(*addBlockCmdPara)
	//		}
	//	}
	case "printchain":
		err := printCmd.Parse(os.Args[2:])
		CheckErr("printChain", err)
		if printCmd.Parsed() {
			cli.PrintChain()
		} else {
			fmt.Println("Not Parsed!!!!!!!!!!!!")
		}
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		CheckErr("getbalanceCmd ", err)
		if getbalanceCmd.Parsed() {
			if *getbalanceCmdPara == "" {
				fmt.Println(Usage)
				os.Exit(1)
			}
			cli.GetBalance(*getbalanceCmdPara)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		CheckErr("sendCmd ", err)
		if sendCmd.Parsed() {
			if *fromPara == "" || *toPara == "" || *amountPara == 0 {
				fmt.Println(Usage)
				os.Exit(1)
			}
			cli.Send(*fromPara, *toPara, *amountPara)

		}

	default:
		fmt.Println("invalid cmd\n", Usage)
		os.Exit(1)
	}
}
