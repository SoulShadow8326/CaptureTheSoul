package controllers
import(
	"bytes"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
)
func StartChallengeContainer(flag string)(int, string, error){
	port := rand.Intn(9999-9000) + 9000
	cmd := exec.Command("docker", "run",
	"-d",
	"-p", fmt.Sprintf("%d:1337", port),
	"-e", fmt.Sprintf("FLAG=%s", flag),
	"vuln-challenge",
)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil{
		return 0, "", fmt.Errorf("docker run failed: %v - %s", err, out.String())
	}
	containerID := strings.TrimSpace(out.String())
	return port, containerID, nil
}	