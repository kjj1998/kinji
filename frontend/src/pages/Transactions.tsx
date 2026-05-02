import { Button, Group, Select, Text, TextInput } from "@mantine/core";
import "@mantine/dates/styles.css";
import type { ColDef } from "ag-grid-community";
import { AllCommunityModule } from "ag-grid-community";
import {
	AgGridProvider,
	AgGridReact,
	type CustomCellRendererProps,
} from "ag-grid-react";
import { useMemo, useRef, useState } from "react";
import { AddTransactionForm } from "../components";
import { transactions } from "../data/transactions";
import type { Transaction } from "../types";

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

const CATEGORIES = [
	"Food",
	"Transport",
	"Groceries",
	"Subscriptions",
	"Income",
	"Shopping",
	"Health",
	"Utilities",
];

export function Transactions() {
	const [allTransactions, setAllTransactions] =
		useState<Transaction[]>(transactions);
	const [rowData, setRowData] = useState<Transaction[]>(transactions);
	const [showForm, setShowForm] = useState(false);

	const colDefs = useMemo<ColDef<Transaction>[]>(
		() => [
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
				editable: true,
				cellEditor: "agNumberCellEditor",
				cellEditorParams: { precision: 2 },
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
		],
		[],
	);

	const gridRef = useRef<AgGridReact>(null);

	return (
		<>
			<Group h={60} px="md" style={{ borderBottom: "1px solid #D4A853" }}>
				<Text fw={400} size="xl">
					Transactions
				</Text>
			</Group>
			<Group justify="space-between" mt="md" mb="sm">
				<Group>
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
						data={CATEGORIES}
						onChange={(value) => {
							setRowData(
								value
									? allTransactions.filter((t) => t.category === value)
									: allTransactions,
							);
						}}
					/>
				</Group>
				<Button variant="filled" onClick={() => setShowForm((prev) => !prev)}>
					{showForm ? "Cancel" : "+ Add Transaction"}
				</Button>
			</Group>
			{showForm && (
				<AddTransactionForm
					categories={CATEGORIES}
					onAdd={(newTx) => {
						setAllTransactions((prev) => [newTx, ...prev]);
						setShowForm(false);
					}}
				/>
			)}
			<AgGridProvider modules={modules}>
				<div
					style={{
						height: showForm ? "calc(100vh - 365px)" : "calc(100vh - 180px)",
					}}
				>
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
