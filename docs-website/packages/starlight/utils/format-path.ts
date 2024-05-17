import project from 'virtual:starlight/project-context';
import { createPathFormatter } from './createPathFormatter';

export const formatPath = createPathFormatter({
	format: project.build.format,
	trailingSlash: project.trailingSlash,
});
