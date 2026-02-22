import { toast } from 'svelte-sonner';

export function showSuccess(message: string) {
	toast.success(message);
}

export function showError(error: any, fallbackMessage: string = 'An error occurred') {
	const message = error?.message ?? fallbackMessage;
	toast.error(message);
}

export function showInfo(message: string) {
	toast.info(message);
}

export function showWarning(message: string) {
	toast.warning(message);
}

export const toastMessages = {
	saveSuccess: (itemName: string = 'Changes') => `${itemName} saved successfully.`,
	saveError: (itemName: string = 'changes') => `Failed to save ${itemName}.`,
	deleteSuccess: (itemName: string) => `${itemName} deleted successfully.`,
	deleteError: (itemName: string) => `Failed to delete ${itemName}.`,
	createSuccess: (itemName: string) => `${itemName} created successfully.`,
	createError: (itemName: string) => `Failed to create ${itemName}.`,
	updateSuccess: (itemName: string) => `${itemName} updated successfully.`,
	updateError: (itemName: string) => `Failed to update ${itemName}.`,
	requiredFields: 'Please fill in all required fields.',
	invalidUrl: 'Please enter a valid URL.'
};
