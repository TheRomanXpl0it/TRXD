import os

challs = os.listdir('.')
remotes = [c for c in challs if os.path.isdir(c) and 'remote' in os.listdir(c)]
for r in remotes:
	os.system(f'docker build ./{r}/remote/ -t cc-{r}')
