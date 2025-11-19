package services

import (
  // "database/sql"
  "context"
  "time"
  "fmt"
  "log"
  "sync"

  "github.com/opussocialcontent/opus-go/actions"
  "github.com/opussocialcontent/opus-go/quality"
)

type PubSubHub struct {
// PubSubHub implements a publish-subscribe messaging system
  mu             sync.RWMutex
  topics         map[string]*Topic
  container actions.Container
}

// Topic represents a message channel with subscribers
type Topic struct {
  name       string
  subscribers map[string]HubSubscriber
  mu         sync.RWMutex
}

// Subscriber defines a message handler
type HubSubscriber struct {
  id      string
  handler actions.ActionFunc
  options SubscriberOptions
}

// SubscriberOptions configures subscriber behavior
type SubscriberOptions struct {
  Concurrency int // Number of concurrent message processors
  RetryPolicy RetryPolicy
}

// RetryPolicy defines message retry behavior
type RetryPolicy struct {
  MaxAttempts int
  Backoff     time.Duration
}

// Message wraps payload with metadata
type Message struct {
  Topic   string
  Payload actions.Payload
  Context context.Context
}

// NewPubSubHub creates a new PubSubHub instance
func NewPubSubHub(container actions.Container) *PubSubHub {
  return &PubSubHub{
    topics:         make(map[string]*Topic),
    container: container,
  }
}

// CreateTopic initializes a new topic
func (a *PubSubHub) CreateTopic(name string) error {
  a.mu.Lock()
  defer a.mu.Unlock()

  if _, exists := a.topics[name]; exists {
    return quality.ErrTopicExists
  }

  a.topics[name] = &Topic{
    name:       name,
    subscribers: make(map[string]HubSubscriber),
  }
  return nil
}

// Subscribe registers an action handler for a topic
func (a *PubSubHub) Subscribe(topicName, actionName string, opts SubscriberOptions) error {
  a.mu.RLock()
  topic, exists := a.topics[topicName]
  a.mu.RUnlock()
  if !exists {
    return quality.ErrTopicNotFound
  }

  // Resolve action from registry
  registry := a.container.Registry()
  actionFunc, ok := registry.ResolveAction(actionName)
  if !ok {
    return quality.ErrInvalidAction.WithDetail(fmt.Sprintf("action %s not found", actionName))
  }

  topic.mu.Lock()
  defer topic.mu.Unlock()

  if _, exists := topic.subscribers[actionName]; exists {
    return quality.ErrSubscriberExists
  }

  topic.subscribers[actionName] = HubSubscriber{
    id:      actionName,
    handler: actionFunc,
    options: opts,
  }

  return nil
}
// Publish sends a message to all subscribers of a topic
func (a *PubSubHub) Publish(topicName string, p actions.Payload) error {
  fmt.Println(topicName)
  fmt.Println(topicName)
  fmt.Println(topicName)
  fmt.Println(topicName)
  a.mu.RLock()
  topic, exists := a.topics[topicName]
  if !exists {
    a.mu.RUnlock()
    return quality.ErrTopicNotFound
  }
  fmt.Println(topic)
  fmt.Println(topic)
  fmt.Println(topic)
  fmt.Println(topic)
  fmt.Println(topic)
  
  // Get a copy of subscribers while holding the lock to avoid race conditions
  topic.mu.RLock()
  subscribers := make([]HubSubscriber, 0, len(topic.subscribers))
  for _, sub := range topic.subscribers {
    fmt.Println("subscribers loop")
    subscribers = append(subscribers, sub)
  }
  fmt.Println(len(subscribers))
  fmt.Println(len(subscribers))
  fmt.Println(len(subscribers))
  fmt.Println(len(subscribers))
  fmt.Println(len(subscribers))
  fmt.Println(len(subscribers))
  fmt.Println(len(subscribers))
  fmt.Println(len(subscribers))
  fmt.Println(len(subscribers))
  fmt.Println(len(subscribers))
  fmt.Println(len(subscribers))
  fmt.Println(len(subscribers))
  topic.mu.RUnlock()
  a.mu.RUnlock() // Release the main lock after copying subscribers

  msg := Message{
    Topic:   topicName,
    Payload: p,
    Context: context.Background(), // Add context if needed
  }

  // Process messages concurrently for each subscriber
  var wg sync.WaitGroup
  for _, sub := range subscribers {
    wg.Add(1)
    go func(s HubSubscriber) {
      defer wg.Done()
      a.processMessage(msg, s)
    }(sub)
  }

  wg.Wait()
  return nil
}

// processMessage handles message delivery with retry logic
func (a *PubSubHub) processMessage(msg Message, sub HubSubscriber) {
  fmt.Println("PROCESSING")
  fmt.Println("PROCESSING")
  fmt.Println("PROCESSING")
  fmt.Println("PROCESSING")
  fmt.Println("PROCESSING")
  var lastErr error
  opts := sub.options

  for attempt := 1; attempt <= opts.RetryPolicy.MaxAttempts; attempt++ {
    // Create action with the resolved handler
    action := actions.NewAction(
      fmt.Sprintf("pubsub:%s:%s", msg.Topic, sub.id),
      msg.Payload,
      a.container,
      sub.handler,
    )

    if err := action.Execute(); err != nil {
      lastErr = err
      log.Printf("Attempt %d failed for subscriber %s: %v", attempt, sub.id, err)
      if attempt < opts.RetryPolicy.MaxAttempts {
        time.Sleep(opts.RetryPolicy.Backoff * time.Second * time.Duration(attempt))
      }
      continue
    }

    return
  }

  log.Printf("Message processing failed after %d attempts for subscriber %s: %v", 
    opts.RetryPolicy.MaxAttempts, sub.id, lastErr)
}

// Shutdown gracefully stops all message processing
func (a *PubSubHub) Shutdown() {
  a.mu.Lock()
  defer a.mu.Unlock()

  // Clear all topics and subscribers
  a.topics = make(map[string]*Topic)
}


