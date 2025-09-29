Actions & Payloads
An action is an implementation of the Command Pattern (). When an action is created, it receives a data payload and a container with service depencies. When its executed, it maps data to or from the service in container.

Consumers & Providers
A provider is a set of actions and payloads contained in a domain. A consumer is a service that calls actions. The consumer does not need to match the provider. Any consumer can call any number of actions from any provider.
Consumers can:
 - route actions to api
 - route actions to views
 - route static views
 - register actions as cron jobs
 - subscribe actions to pubsub events


