import "@mantine/dates/styles.css";
import { useQueryClient } from "@tanstack/react-query";
import type { ColDef } from "ag-grid-community";
import { AllCommunityModule } from "ag-grid-community";
import {
	AgGridProvider,
	AgGridReact,
	type CustomCellRendererProps,
} from "ag-grid-react";
import { useMemo, useRef, useState } from "react";
import { AddTransactionForm, Header, TransactionFilters } from "../components";
import { categories } from "../data";
import { useTransactions } from "../hooks";
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

export function Transactions() {
	const queryClient = useQueryClient();
	const { data: allTransactions = [] } = useTransactions("james");
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

	console.log(allTransactions);

	return (
		<>
			<Header text={"Transactions"} />
			<TransactionFilters
				categories={categories}
				showForm={showForm}
				onTextInputChange={(e) =>
					gridRef.current?.api.setGridOption(
						"quickFilterText",
						e.currentTarget.value,
					)
				}
				onCategorySelectChange={(value) => {
					gridRef.current?.api
						.setColumnFilterModel(
							"category",
							value
								? { filterType: "text", type: "equals", filter: value }
								: null,
						)
						.then(() => gridRef.current?.api.onFilterChanged());
				}}
				onShowFormClick={() => setShowForm((prev) => !prev)}
			/>
			{showForm && (
				<AddTransactionForm
					categories={categories}
					onAdd={(newTx) => {
						queryClient.setQueryData(
							["transactions", "james"],
							(prev: Transaction[]) => [newTx, ...prev],
						);
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
						rowData={allTransactions}
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
