Nobody likes flaky tests, but they still occur. There are many reasons for this, such as depending on the order of keys in a map or using `time.Sleep` to wait for something to happen.

## Flaky Test

For example, you need to start a server that could take a couple of seconds (1-5 seconds) to start, and you want to test it:
<details>
<summary>Yikes tests that takes more than 1 second!</summary>
Sometimes we need it, but it's not a good practice when tests take a lot of time to complete. You must have a fast feedback loop.
</details>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ go
func startServer() {
	rand.Seed(time.Now().UnixNano())
	sleepDuration := time.Duration(rand.Intn(5)+1) * time.Second
	time.Sleep(sleepDuration)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	http.ListenAndServe(":8080", nil)
}
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Then you write a flaky test with `time.Sleep` for it:
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ go
func TestServer(t *testing.T) {
  go startServer()
  time.Sleep(1 * time.Second)
  resp, err := http.Get("http://localhost:8080")
  if err != nil {
    t.Fatal(err)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    t.Fatal(err)
  }
  if string(body) != "Hello, world!" {
    t.Fatalf("unexpected body: %s", body)
  }
}
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~


What is the problem with this?
It's flaky; sometimes it will fail, sometimes it will pass.
There's no guarantee that the server will start in 1 second.
But how can we fix this?

## Using Awaitility

We can use [Awaitility](https://github.com/mehXX/awaitility) to wait for the server to start
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ go
func TestServer(t *testing.T) {
    go startServer()

    ctx := context.Background()
    err := awaitility.Await(ctx, 100*time.Millisecond, 5*time.Second, func() bool {
        resp, err := http.Get("http://localhost:8080/")
        if err != nil {
            return false
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return false
        }

        return string(body) == "Hello, world!"
    })

    if err != nil {
        t.Errorf("Unexpected error during await: %s", err)
    }
}

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Instead of starting server there could be anything like publishing a message to a queue or waiting that something would happen in another goroutine.

you can find full code here:

awaitility test: [https://github.com/holdmycatnip/mehXX.github.io/blob/master/_posts/you-must-avoid-sleep-at-all-cost/awaitility_test.go](https://github.com/mehXX/mehXX.github.io/blob/master/_posts/you-must-avoid-sleep-at-all-cost/awaitility_test.go)

flaky test: [https://github.com/holdmycatnip/mehXX.github.io/blob/master/_posts/you-must-avoid-sleep-at-all-cost/flaky_test.go](https://github.com/mehXX/mehXX.github.io/blob/master/_posts/you-must-avoid-sleep-at-all-cost/flaky_test.go)
