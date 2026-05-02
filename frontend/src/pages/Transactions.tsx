import { Group, Select, Text, TextInput } from "@mantine/core";
import type { ColDef } from "ag-grid-community";
import { AllCommunityModule } from "ag-grid-community";
import {
	AgGridProvider,
	AgGridReact,
	type CustomCellRendererProps,
} from "ag-grid-react";
import { useRef, useState } from "react";
import { type Transaction, transactions } from "../data/transactions";

const modules = [AllCommunityModule];

function AmountCell(props: CustomCellRendererProps<Transaction, number>) {
	const value = props.value ?? 0;
	let displayValue = "";

	if (value > 0) {
		displayValue = `+$${value.toFixed(2)}`;
	} else if (value < 0) {
		displayValue = `-$${Math.abs(value).toFixed(2)}`;
	} else {
		displayValue = "$0.00";
	}
	return (
		<span style={{ color: value > 0 ? "green" : "inherit" }}>
			{displayValue}
		</span>
	);
}

function SplitCell(props: CustomCellRendererProps<Transaction, number>) {
	const value = props.value;

	if (value == null) {
		return null;
	}

	return <span style={{ color: "grey" }}>${Math.abs(value).toFixed(2)}</span>;
}

export function Transactions() {
	const [allTransactions] = useState<Transaction[]>(transactions);
	const [rowData, setRowData] = useState<Transaction[]>(transactions);

	const [colDefs, _setColDefs] = useState<ColDef<Transaction>[]>([
		{ field: "date", headerName: "Date", sortable: true, flex: 1 },
		{
			field: "merchant",
			headerName: "Merchant",
			filter: "agTextColumnFilter",
			flex: 2,
		},
		{
			field: "category",
			headerName: "Category",
			filter: "agTextColumnFilter",
			flex: 1,
		},
		{
			field: "amount",
			headerName: "Amount",
			sortable: true,
			cellRenderer: AmountCell,
			flex: 1,
		},
		{
			field: "notes",
			headerName: "Notes",
			editable: true,
			cellEditor: "agTextCellEditor",
			flex: 3,
		},
		{
			field: "split",
			headerName: "Split",
			cellRenderer: SplitCell,
			editable: true,
			cellEditor: "agNumberCellEditor",
			cellEditorParams: { min: 0, precision: 2 },
			flex: 1,
		},
	]);

	const gridRef = useRef<AgGridReact>(null);

	return (
		<>
			<Group h={60} px="md" style={{ borderBottom: "1px solid #D4A853" }}>
				<Text fw={400} size="xl">
					Transactions
				</Text>
			</Group>
			<Group mt="md" mb="sm">
				<TextInput
					placeholder="Search..."
					style={{ width: 240 }}
					onChange={(e) =>
						gridRef.current?.api.setGridOption(
							"quickFilterText",
							e.currentTarget.value,
						)
					}
				/>
				<Select
					placeholder="All Categories"
					clearable
					style={{ width: 180 }}
					maxDropdownHeight={400}
					data={[
						"Food",
						"Transport",
						"Groceries",
						"Subscriptions",
						"Income",
						"Shopping",
						"Health",
						"Utilities",
					]}
					onChange={(value) => {
						setRowData(
							value
								? allTransactions.filter((t) => t.category === value)
								: allTransactions,
						);
					}}
				/>
			</Group>
			<AgGridProvider modules={modules}>
				<div style={{ height: "calc(100vh - 180px)" }}>
					<AgGridReact
						rowData={rowData}
						columnDefs={colDefs}
						pagination={true}
						paginationPageSize={20}
						ref={gridRef}
						defaultColDef={{ flex: 1 }}
					/>
				</div>
			</AgGridProvider>
		</>
	);
}
