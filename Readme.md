Things to improve:

1. Handle context properly
2. Use state machine for files states.This will help to resolve the following issues:

- Multiple boolean columns in DB, that potentially can be resolved only by one column that shows the current state of the file
- RemoveUploadedFiles crontask runs continuously for the handled files. So it makes redundant calls to DB and uses additional resources

3. Better resolving servers addresses. Maybe use Redis as a common address resolver. Need to research this topic.
4. Add graceful shutdown
5. Better error handling
6. Better folders structure
7. Use transactions
