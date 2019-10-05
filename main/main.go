package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	id                   int
	isAlive              bool
	term                 int
	voted                bool
	votedFor             int
	electionTimeout      int
	electionTimeoutClock int
	state                string
	log                  []string

	commitIndex int
	lastApplied int
	nextIndex   []int
	matchIndex  []int

	res int
}

// func AppendEntries(nodes []*Node, leaderIndex int, term int, prevLogIndex int, prevLogTerm int, entries []string, leaderCommitIndex int) {
func AppendEntries(nodes []*Node, leaderIndex int, entries []string) {
	// prevLogIndex := nodes[leaderIndex].lastApplied
	// prevLogTerm := nodes[leaderIndex].term
	// leaderCommitIndex := nodes[leaderIndex].commitIndex

	for _, node := range nodes {
		if node.isAlive == true {
			node.electionTimeoutClock = node.electionTimeout
		}
	}
}

// func RequestVote(term int, candidateId int, lastLogIndex int, lastLogTerm int){
// 	//
// }

func InitializeNodes(n int) []*Node {
	var nodes []*Node
	indices := make([]int, n, n)
	for i := 0; i < n; i++ {
		indices[i] = 1
	}
	for i := 0; i < n; i++ {
		nodes = append(nodes, &Node{id: i, state: "follower", term: 0, voted: false, commitIndex: 0, lastApplied: 0, nextIndex: indices, res: 0, isAlive: true})
	}
	return nodes
}

func setElectionTimeout(node *Node, minElectionTimeout int, maxElectionTimeout int) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	node.state = "follower"
	node.electionTimeout = r.Intn(maxElectionTimeout-minElectionTimeout+1) + minElectionTimeout
	node.electionTimeoutClock = node.electionTimeout
}

func reduceTimeout(nodes []*Node) []int {
	var candidateId []int
	for i, node := range nodes {
		node.electionTimeoutClock--
		if node.isAlive == false {
			node.electionTimeoutClock = node.electionTimeout
		}
		if node.electionTimeoutClock == 0 && node.isAlive {
			node.state = "candidate"
			node.term++
			fmt.Println("A candidate arises", i, "\n")
			candidateId = append(candidateId, i) //or append(candidateId, node.id) ;id and index 'i' are same here
		}
	}
	return candidateId
}

func show(nodes []*Node) {
	fmt.Println("ID isAlive term voted votedFor electionTimeout electionTimeoutClock state log")
	for _, node := range nodes {
		fmt.Printf("%v\n", node)
	}
}

func saveMessage(msg string, nodes []*Node) {
	for _, node := range nodes {
		node.log = append(node.log, msg)
	}
}

func runElection(nodes []*Node, candidateId []int) int {
	fmt.Println("\n...RUNNING ELECTION...\n")
	// cnt := 0
	voteCounter := make([]int, len(candidateId))

	//RUNNING OUT OF MEMORY HERE

	//Cast vote
	//vote for oneself
	for i, ind := range candidateId {
		nodes[ind].voted = true
		nodes[ind].votedFor = ind
		voteCounter[i]++
		fmt.Println("Node", ind, "voted for Candidate", ind)
	}

	//followers' vote
	for n, node := range nodes {
		//
		for i, candidate_id := range candidateId {
			if node.isAlive && node.voted != true {
				if nodes[candidate_id].term >= node.term {
					node.voted = true
					node.votedFor = i
					voteCounter[i]++
					fmt.Println("Node", n, "voted for Candidate", candidate_id)
					break
				}
			}
		}

		// //vote for oneself
		// if cnt < len(candidateIndex) && i == candidateIndex[cnt] {
		// 	node.voted = true
		// 	node.votedFor = i
		// 	voteCounter[cnt]++
		// 	fmt.Println("Node", i, "voted for Candidate", candidateIndex[cnt])

		// 	cnt++
		// } else {
		// 	// only vote if node is alive
		// 	if nodes[i].isAlive {
		// 		nodes[i].term++
		// 		//current implementation: random vote
		// 		//work remaining: vote based on criteria set by paper
		// 		r := rand.NewSource(time.Now().UnixNano())
		// 		s := rand.New(r)
		// 		random_vote := s.Intn(len(candidateIndex))
		// 		voteCounter[random_vote]++
		// 		fmt.Println("Node", i, "voted for Candidate", candidateIndex[random_vote])
		// 	}
		// }
	}

	fmt.Println("...END OF ELECTION...")

	//Calculate maximum vote
	maxVote := 0
	maxVoteCount := 1

	for i := 1; i < len(voteCounter); i++ {
		if voteCounter[i] == voteCounter[maxVote] {
			maxVoteCount++
		}
		if voteCounter[i] > voteCounter[maxVote] {
			maxVote = i
			maxVoteCount = 1
		}
	}

	//return -1 (error) if 2 or more candidates have equal max votes
	if maxVoteCount > 1 {
		return -1
	}

	return candidateId[maxVote]
}

func main() {
	done := false                        //to run the loop until stopped
	cnt := 0                             //logical clock
	hearbeatInterval := 1                //interval at which heartbeat msg is sent by the leader
	heartbeatCounter := hearbeatInterval //monitor if it's time to send heartbeat msg
	var leaderId int = -1                //leader's position in the array of *Node
	maxElectionTimeout := 8
	minElectionTimeout := 4
	nodes := InitializeNodes(5) //array of *Node

	for _, node := range nodes {
		setElectionTimeout(node, minElectionTimeout, maxElectionTimeout) //initialize election timeout value
	}

	for done != true {
		heartbeatCounter--
		fmt.Println("---------------------------")
		fmt.Printf("Iteration %d :", cnt+1)
		reader := bufio.NewReader(os.Stdin)
		byte_val, _, _ := reader.ReadLine() //read AppendEntry command
		string_val := string(byte_val)
		msg := strings.Split(string_val, " ")
		msg_handled := false

		// Introduce server crash and recovery
		switch msg[0] {
		case "crash":
			msg_handled = true
			node_id, _ := strconv.Atoi(msg[1])
			nodes[node_id].isAlive = false
		case "recover":
			msg_handled = true
			node_id, _ := strconv.Atoi(msg[1])
			nodes[node_id].isAlive = true

			//also look if the recovered system was a leader before it crashed
			//handle the situation if it was a leader
			if nodes[node_id].state == "leader" {
				//todo: follower, candidate or leader check here
				if nodes[node_id].term > nodes[leaderId].term {
					leaderId = node_id
				} else {
					nodes[node_id].state = "follower"
				}
			}
		}

		if !msg_handled && string_val != "6" && string_val != "" {
			saveMessage(string_val, nodes) //save the command (to be committed after leader confirms commit)
		}

		//handle edge case when heartbeatInterval is set to 0
		if heartbeatCounter < 0 {
			heartbeatCounter = 0
		}
		/*
		   if len(candidateIndex)==1, leader = candidateIndex
		   if len(candidateIndex)>1, resetElectionTimeout(candidateIndex)
		*/

		if heartbeatCounter == 0 && leaderId != -1 && nodes[leaderId].isAlive == true { //if leader exists and hasn't crashed
			//send heartbeat message
			AppendEntries(nodes, leaderId, msg)
		} else { //No heartbeat message: no leader or not time for message sending yet
			candidateId := reduceTimeout(nodes)
			if len(candidateId) != 0 {
				show(nodes) //show node state before leader voting

				//todo: even if 1 candidate, ask for vote
				//fix after adding voting criteria

				fmt.Println("Candidates", candidateId)

				leaderId = runElection(nodes, candidateId)

				if leaderId == -1 { //leader appointment failed
					for _, nodeId := range candidateId {
						setElectionTimeout(nodes[nodeId], minElectionTimeout, maxElectionTimeout)
					}
				} else { //successful leader appointment
					AppendEntries(nodes, leaderId, msg)

					nodes[leaderId].state = "leader"
					fmt.Println("\nNew leader selected: ", leaderId)

					//demote all candidates to followers
					for _, ind := range candidateId {
						if nodes[ind].state == "candidate" && ind != leaderId {
							nodes[ind].state = "follower"
						}
					}
				}

			}
		}
		if heartbeatCounter == 0 {
			heartbeatCounter = hearbeatInterval
		}

		switch string_val {
		case "6": //terminate if "6" is pressed
			done = true
			break
		case "":
			show(nodes)
		default:
			show(nodes)
		}

		cnt++
		// done = true
	}
}
