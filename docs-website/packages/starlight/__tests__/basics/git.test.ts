import { mkdtempSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { tmpdir } from 'node:os';
import { spawnSync } from 'node:child_process';
import { describe, expect, test } from 'vitest';
import { getNewestCommitDate } from '../../utils/git';

describe('getNewestCommitDate', () => {
	const { commitAllChanges, getFilePath, writeFile } = makeTestRepo();

	test('returns the newest commit date', () => {
		const file = 'updated.md';
		const lastCommitDate = '2023-06-25';

		writeFile(file, 'content 0');
		commitAllChanges('add updated.md', '2023-06-21');
		writeFile(file, 'content 1');
		commitAllChanges('update updated.md', lastCommitDate);

		expectCommitDateToEqual(getNewestCommitDate(getFilePath(file)), lastCommitDate);
	});

	test('returns the initial commit date for a file never updated', () => {
		const file = 'added.md';
		const commitDate = '2022-09-18';

		writeFile(file, 'content');
		commitAllChanges('add added.md', commitDate);

		expectCommitDateToEqual(getNewestCommitDate(getFilePath(file)), commitDate);
	});

	test('returns the newest commit date for a file with a name that contains a space', () => {
		const file = 'updated with space.md';
		const lastCommitDate = '2021-01-02';

		writeFile(file, 'content 0');
		commitAllChanges('add updated.md', '2021-01-01');
		writeFile(file, 'content 1');
		commitAllChanges('update updated.md', lastCommitDate);

		expectCommitDateToEqual(getNewestCommitDate(getFilePath(file)), lastCommitDate);
	});

	test('returns the newest commit date for a file updated the same day', () => {
		const file = 'updated-same-day.md';
		const lastCommitDate = '2023-06-25T14:22:35Z';

		writeFile(file, 'content 0');
		commitAllChanges('add updated.md', '2023-06-25T12:34:56Z');
		writeFile(file, 'content 1');
		commitAllChanges('update updated.md', lastCommitDate);

		expectCommitDateToEqual(getNewestCommitDate(getFilePath(file)), lastCommitDate);
	});

	test('throws when failing to retrieve the git history for a file', () => {
		expect(() => getNewestCommitDate(getFilePath('../not-a-starlight-test-repo/test.md'))).toThrow(
			/^Failed to retrieve the git history for file "[/\\:-\w ]+[/\\]test\.md"/
		);
	});

	test('throws when trying to get the history of a non-existing or untracked file', () => {
		const expectedError =
			/^Failed to validate the timestamp for file "[/\\:-\w ]+[/\\](?:unknown|untracked)\.md"$/;
		writeFile('untracked.md', 'content');

		expect(() => getNewestCommitDate(getFilePath('unknown.md'))).toThrow(expectedError);
		expect(() => getNewestCommitDate(getFilePath('untracked.md'))).toThrow(expectedError);
	});
});

function expectCommitDateToEqual(commitDate: CommitDate, expectedDateStr: ISODate) {
	const expectedDate = new Date(expectedDateStr);
	expect(commitDate).toStrictEqual(expectedDate);
}

function makeTestRepo() {
	const repoPath = mkdtempSync(join(tmpdir(), 'starlight-test-git-'));

	function runInRepo(command: string, args: string[], env: NodeJS.ProcessEnv = {}) {
		const result = spawnSync(command, args, { cwd: repoPath, env });

		if (result.status !== 0) {
			throw new Error(`Failed to execute test repository command: '${command} ${args.join(' ')}'`);
		}
	}

	// Configure git specifically for this test repository.
	runInRepo('git', ['init']);
	runInRepo('git', ['config', 'user.name', 'starlight-test']);
	runInRepo('git', ['config', 'user.email', 'starlight-test@example.com']);
	runInRepo('git', ['config', 'commit.gpgsign', 'false']);

	return {
		// The `dateStr` argument should be in the `YYYY-MM-DD` or `YYYY-MM-DDTHH:MM:SSZ` format.
		commitAllChanges(message: string, dateStr: ISODate) {
			const date = dateStr.endsWith('Z') ? dateStr : `${dateStr}T00:00:00Z`;

			runInRepo('git', ['add', '-A']);
			// This sets both the author and committer dates to the provided date.
			runInRepo('git', ['commit', '-m', message, '--date', date], { GIT_COMMITTER_DATE: date });
		},
		getFilePath(name: string) {
			return join(repoPath, name);
		},
		writeFile(name: string, content: string) {
			writeFileSync(join(repoPath, name), content);
		},
	};
}

type ISODate =
	| `${number}-${number}-${number}`
	| `${number}-${number}-${number}T${number}:${number}:${number}Z`;

type CommitDate = ReturnType<typeof getNewestCommitDate>;
