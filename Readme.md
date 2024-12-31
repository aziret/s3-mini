# How to start
Run:
`
docker compose up
`

This will create a main frontend and backend on the dedicated ports, also it will create 6 separate storages to store file chunks in them.

# Description
The app mainly uses crontasks to:
- create file chunks for the given file
- upload file chunks to the provided storages
- remove files from the main server
- when requested, it will compose a single file from the chunks stored in the storages and give it to the user

## Things to improve:

1. Handle context properly
2. Use state machine for files states.This will help to resolve the following issues:

- Multiple boolean columns in DB, that potentially can be resolved only by one column that shows the current state of the file
- RemoveUploadedFiles crontask runs continuously for the handled files. So it makes redundant calls to DB and uses additional resources

3. Better resolving servers addresses. Maybe use Redis as a common address resolver. Need to research this topic.
4. Add graceful shutdown
5. Better error handling
6. Better folders structure
7. Use transactions
8. Use cache for downloaded files and remove the downloaded file from the main server
