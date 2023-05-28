.PHONY: create-tigerbeetle-db-file tigerbeetle temporalite run clean

run-tigerbeetle: clean create-tigerbeetle-db-file tigerbeetle

create-tigerbeetle-db-file:
	tigerbeetle format --cluster=0 --replica=0 --replica-count=1 0_0.tigerbeetle

tigerbeetle:
	tigerbeetle start --addresses=3000 0_0.tigerbeetle

temporalite:
	temporalite start --namespace default --ephemeral

run:
	encore run

clean:
	rm 0_0.tigerbeetle
