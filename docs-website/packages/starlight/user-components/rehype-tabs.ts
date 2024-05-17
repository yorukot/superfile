import type { Element } from 'hast';
import { select } from 'hast-util-select';
import { rehype } from 'rehype';
import { CONTINUE, SKIP, visit } from 'unist-util-visit';
import { Icons } from '../components/Icons';

interface Panel {
	panelId: string;
	tabId: string;
	label: string;
	icon?: keyof typeof Icons;
}

declare module 'vfile' {
	interface DataMap {
		panels: Panel[];
	}
}

export const TabItemTagname = 'starlight-tab-item';

// https://github.com/adobe/react-spectrum/blob/99ca82e87ba2d7fdd54f5b49326fd242320b4b51/packages/%40react-aria/focus/src/FocusScope.tsx#L256-L275
const focusableElementSelectors = [
	'input:not([disabled]):not([type=hidden])',
	'select:not([disabled])',
	'textarea:not([disabled])',
	'button:not([disabled])',
	'a[href]',
	'area[href]',
	'summary',
	'iframe',
	'object',
	'embed',
	'audio[controls]',
	'video[controls]',
	'[contenteditable]',
	'[tabindex]:not([disabled])',
]
	.map((selector) => `${selector}:not([hidden]):not([tabindex="-1"])`)
	.join(',');

let count = 0;
const getIDs = () => {
	const id = count++;
	return { panelId: 'tab-panel-' + id, tabId: 'tab-' + id };
};

/**
 * Rehype processor to extract tab panel data and turn each
 * `<starlight-tab-item>` into a `<section>` with the necessary
 * attributes.
 */
const tabsProcessor = rehype()
	.data('settings', { fragment: true })
	.use(function tabs() {
		return (tree: Element, file) => {
			file.data.panels = [];
			let isFirst = true;
			visit(tree, 'element', (node) => {
				if (node.tagName !== TabItemTagname || !node.properties) {
					return CONTINUE;
				}

				const { dataLabel, dataIcon } = node.properties;
				const ids = getIDs();
				const panel: Panel = {
					...ids,
					label: String(dataLabel),
				};
				if (dataIcon) panel.icon = String(dataIcon) as keyof typeof Icons;
				file.data.panels?.push(panel);

				// Remove `<TabItem>` props
				delete node.properties.dataLabel;
				delete node.properties.dataIcon;
				// Turn into `<section>` with required attributes
				node.tagName = 'section';
				node.properties.id = ids.panelId;
				node.properties['aria-labelledby'] = ids.tabId;
				node.properties.role = 'tabpanel';

				const focusableChild = select(focusableElementSelectors, node);
				// If the panel does not contain any focusable elements, include it in
				// the tab sequence of the page.
				if (!focusableChild) {
					node.properties.tabindex = 0;
				}

				// Hide all panels except the first
				// TODO: make initially visible tab configurable
				if (isFirst) {
					isFirst = false;
				} else {
					node.properties.hidden = true;
				}

				// Skip over the tab panel’s children.
				return SKIP;
			});
		};
	});

/**
 * Process tab panel items to extract data for the tab links and format
 * each tab panel correctly.
 * @param html Inner HTML passed to the `<Tabs>` component.
 */
export const processPanels = (html: string) => {
	const file = tabsProcessor.processSync({ value: html });
	return {
		/** Data for each tab panel. */
		panels: file.data.panels,
		/** Processed HTML for the tab panels. */
		html: file.toString(),
	};
};
