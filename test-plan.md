# Test Plan

Sorry, I don't have unit tests for this code.
Instead, I have this test instruction list.
These are a series of actions to try out on a Kubernetes cluster
to make sure the tool is working as expected

## Things to Check For

- Do permissions work as expected?
  - Can the service do what it's supposed to be able to do?
  - Is the service blocked when it should be blocked?
- Check different combinations of namespace values
  - None
  - One constant
  - Multiple constant
  - Constants and wildcards

## Checking Permissions

Permissions are on a namespace level. Create the following test namespaces:

- test-1: Has permission
- test-2: Does not have permission
- test-3: Has permission
- test-11: Has permission
- test-12: Does not have permission
- test-13: Has permission

Try the following values for `TARGET_NAMESAPCE`:

- None (empty)
- `*`
- `test-1`
- `test-1,test-2,test-3`
- `test-?`
- `test-1?`
- `test-*`
 