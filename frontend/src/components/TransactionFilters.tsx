import { Button, Group, Select, TextInput } from "@mantine/core";

const MONTHS = [
	{ value: "01", label: "January" },
	{ value: "02", label: "February" },
	{ value: "03", label: "March" },
	{ value: "04", label: "April" },
	{ value: "05", label: "May" },
	{ value: "06", label: "June" },
	{ value: "07", label: "July" },
	{ value: "08", label: "August" },
	{ value: "09", label: "September" },
	{ value: "10", label: "October" },
	{ value: "11", label: "November" },
	{ value: "12", label: "December" },
];

export interface TransactionFiltersProps {
	categories: string[];
	years: string[];
	availableMonths: string[];
	selectedMonth: string | null;
	selectedYear: string | null;
	showForm: boolean;
	onTextInputChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
	onCategorySelectChange: (value: string | null) => void;
	onMonthSelectChange: (value: string | null) => void;
	onYearSelectChange: (value: string | null) => void;
	onShowFormClick: () => void;
}

export function TransactionFilters({
	categories,
	years,
	availableMonths,
	selectedMonth,
	selectedYear,
	showForm,
	onTextInputChange,
	onCategorySelectChange,
	onMonthSelectChange,
	onYearSelectChange,
	onShowFormClick,
}: TransactionFiltersProps) {
	const monthOptions =
		availableMonths.length > 0
			? MONTHS.filter((m) => availableMonths.includes(m.value))
			: MONTHS;

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
				<Select
					placeholder="Month"
					clearable
					style={{ width: 140 }}
					data={monthOptions}
					value={selectedMonth}
					onChange={onMonthSelectChange}
				/>
				<Select
					placeholder="Year"
					clearable
					style={{ width: 120 }}
					data={years}
					value={selectedYear}
					onChange={onYearSelectChange}
				/>
			</Group>
			<Button variant="filled" onClick={onShowFormClick}>
				{showForm ? "Cancel" : "+ Add Transaction"}
			</Button>
		</Group>
	);
}
