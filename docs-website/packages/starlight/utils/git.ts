import { basename, dirname } from 'node:path';
import { spawnSync } from 'node:child_process';

export function getNewestCommitDate(file: string) {
	const result = spawnSync('git', ['log', '--format=%ct', '--max-count=1', basename(file)], {
		cwd: dirname(file),
		encoding: 'utf-8',
	});

	if (result.error) {
		throw new Error(`Failed to retrieve the git history for file "${file}"`);
	}
	const output = result.stdout.trim();
	const regex = /^(?<timestamp>\d+)$/;
	const match = output.match(regex);

	if (!match?.groups?.timestamp) {
		throw new Error(`Failed to validate the timestamp for file "${file}"`);
	}

	const timestamp = Number(match.groups.timestamp);
	const date = new Date(timestamp * 1000);
	return date;
}
