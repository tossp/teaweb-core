package agents

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/utils/string"
	"testing"
)

func TestAddManyAgents(t *testing.T) {
	count := 1000
	for i := 0; i < count; i++ {
		agentList, err := agents.SharedAgentList()
		if err != nil {
			t.Fatal(err)
		}

		agent := agents.NewAgentConfig()
		agent.On = false
		agent.Name = "Web" + fmt.Sprintf("%d", i)
		agent.Host = "192.168.0." + fmt.Sprintf("%d", i)
		agent.AllowAll = true
		agent.Allow = []string{}
		agent.Key = stringutil.Rand(32)
		//agent.GroupIds = []string{"2kMMzOcWWPFrhdaM"}
		err = agent.Save()
		if err != nil {
			t.Fatal(err)
		}

		agentList.AddAgent(agent.Filename())
		err = agentList.Save()
		if err != nil {
			t.Fatal(err)
		}
	}
}
