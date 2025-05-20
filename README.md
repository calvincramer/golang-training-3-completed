# Golang Training 3 - Concurrency
In `concurrency.go` there are exercises to complete. Try completing them without looking at the tests.

Run the tests with `go test` or `go test -v`.

Passing implementations will be made available at a later time [here](https://github.com/calvincramer/golang-training-3-completed).

# Install Golang (Linux):
1. Download from https://go.dev/dl/
2. Run:
```sh
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz  # use the file downloaded from step 1
```
3. Make sure this path is setup correctly: `export PATH=$PATH:/usr/local/go/bin`. Add this in `.bashrc` or equivalent for your shell.
4. Start a new terminal or reload the current shell environment
5. Verify go works by running `go version`
