type ChangeListener = () => void;

function createMockEditor(initialValue = '') {
	let value = initialValue;
	let changeListener: ChangeListener | undefined;

	return {
		getValue() {
			return value;
		},
		setValue(nextValue: string) {
			value = nextValue;
			changeListener?.();
		},
		onDidChangeModelContent(listener: ChangeListener) {
			changeListener = listener;
			return {
				dispose() {
					changeListener = undefined;
				}
			};
		},
		dispose() {
			changeListener = undefined;
		}
	};
}

export const editor = {
	create(_container: HTMLElement, options?: { value?: string }) {
		return createMockEditor(options?.value ?? '');
	},
	setTheme(_theme: string) {}
};

const monaco = { editor };

export default monaco;
