# TODO
 - resolve TODOs left into the code
 - convention:
   - change name of chall containers from chall_PORT to chall_TID_CID
   - change name from _TID_CID to _CID_TID
   - change name of challs, network and so on, with $PROJECT_NAME prefix
   - change csrf token name
   - api names more similar (categories_create -> name & categories_delete -> category)
 - tests:
   - tests with removal of trxd-shared (and fault tolerance)
   - discord webhook tests in pipeline
   - tg webhook tests in pipeline
   - add tests with a lot of hash named challs (hard cap for networks on nginx?) (try change IPAM to /30 to preserve addrs)
   - instance expire tests
 - features:
   - https://docs.gofiber.io/api/middleware/helmet/
   - https://doc.traefik.io/traefik/expose/docker/
   - challs SNI traefik
   - N instance limit per team
   - ingress only challenges (verify if useful)
   - challenge remote type (table: (tcp & http) + format: ex. "nc {{host}} {{port}}" & "http://{{host}}:{{port}}")
   - add pagination on "all" GETs
   - invisible teams
   - divide configs by section (secrets, instances, something like this...)
 - utility:
   - configs changable by env
   - better names for binary flags

## Frontend
 - convert .pngs to webp
 - Better tags handling in /challenges page
 - Add custom button themes for first/second/third bloods?
 - Improve challenge creation procedure
 - Fix admin flag submission 
 - ~~Improve "Team join" page on smaller screens~~

## Release
 - config to start and end ctf
 - integration tests (for generic behaviour)
 - tests for distributed functioning
 - role configuraion endpoint
 - dropdown menu container
 - DOCUMENTATION

## Ideas / Features
 - editable homepage & theme
 - instanes page (admin only)
 - deletable instances (from admin)
 - submission page (admin only)
 - deletable submissions (from admin)
 - link to join team
 - endpoint + cli cmd to parse a category/challenge file (maybe yaml) with all infos
 - extract data for ctftime
 - scoreboard freeze (idea: use a table or a view to take a snapshot of the scoreboard, update it every time someone solves if not in freeze time (or just prefetch and compute if not exists))
 - telegram bot for first bloods (and webhook generalization)
 - hash verifier flag or script
 - login via CTFTime
 - default starting points for challenges as global config
 - kube support
 - swarm support
 - ctf stats page
 - flag submit that checks for all challenges and mark solved the one of the correct flag (credits to: midnightsun)
 - writeups for challs into the platform (visible after ctf ends)
 - signed flags (AES ECB or hmac with team id (4 byte) and chall id (4 byte))
 - tls instances: https://github.com/inconshreveable/slt
 - likes and dislikes for challs (only for who actually solved)
 - kick user from team
 - mail server
 - endpoint to store image files (like badges and pfp)
 - flag format validator
