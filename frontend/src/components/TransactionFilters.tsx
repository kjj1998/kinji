import { Button, Group, Select, TextInput } from "@mantine/core";

export interface TransactionFiltersProps {
	categories: string[];
	showForm: boolean;
	onTextInputChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
	onCategorySelectChange: (value: string | null) => void;
	onShowFormClick: () => void;
}

export function TransactionFilters({
	categories,
	showForm,
	onTextInputChange,
	onCategorySelectChange,
	onShowFormClick,
}: TransactionFiltersProps) {
	return (
		<Group justify="space-between" mt="md" mb="sm">
			<Group>
				<TextInput
					placeholder="Search..."
					style={{ width: 240 }}
					onChange={onTextInputChange}
				/>
				<Select
					placeholder="All Categories"
					clearable
					style={{ width: 180 }}
					maxDropdownHeight={400}
					data={categories}
					onChange={onCategorySelectChange}
				/>
			</Group>
			<Button variant="filled" onClick={onShowFormClick}>
				{showForm ? "Cancel" : "+ Add Transaction"}
			</Button>
		</Group>
	);
}
