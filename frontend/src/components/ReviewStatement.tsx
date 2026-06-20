import { Button, Group, Text } from "@mantine/core";
import type { ColDef, ValueFormatterParams } from "ag-grid-community";
import { AllCommunityModule, type GridApi } from "ag-grid-community";
import { AgGridProvider, AgGridReact } from "ag-grid-react";
import { useMemo, useRef } from "react";
import type { Transaction } from "../types";
import { CATEGORIES } from "../types";

const modules = [AllCommunityModule];

function currencyFormatter(params: ValueFormatterParams) {
	return `$${(params.value).toFixed(2)}`;
}

interface ReviewStatementProps {
	transactions: Transaction[];
	save: (rows: Transaction[]) => void;
	cancel: () => void;
}

export function ReviewStatement({
	transactions,
	save,
	cancel,
}: ReviewStatementProps) {
	const colDefs = useMemo<ColDef<Transaction>[]>(
		() => [
			{
				field: "date",
				headerName: "Date",
				sortable: true,
				flex: 1,
				editable: true,
				cellEditor: "agTextCellEditor",
				onCellValueChanged: (params) => {
					const { newValue, oldValue, node } = params;
					if (!node) return;

					const isValid =
						/^\d{4}-\d{2}-\d{2}$/.test(newValue) &&
						!Number.isNaN(Date.parse(newValue));

					if (!isValid) {
						node.setDataValue("date", oldValue); // put the old value back
					}
				},
			},
			{
				field: "merchant",
				headerName: "Merchant",
				filter: "agTextColumnFilter",
				flex: 2,
				editable: true,
			},
			{
				field: "category",
				headerName: "Category",
				filter: "agTextColumnFilter",
				flex: 1,
				cellEditor: "agSelectCellEditor",
				cellEditorParams: {
					values: CATEGORIES,
				},
				editable: true,
			},
			{
				field: "amount",
				headerName: "Amount",
				sortable: true,
				editable: true,
				cellEditor: "agNumberCellEditor",
				cellEditorParams: { precision: 2 },
				flex: 1,
				valueGetter: (p) => (p.data?.amount ?? 0) / 100, // cents → dollars (display + editor)
				valueSetter: (p) => {
					p.data.amount = Math.round(p.newValue * 100);
					return true;
				},
				valueFormatter: currencyFormatter,
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
				editable: true,
				cellEditor: "agNumberCellEditor",
				cellEditorParams: { min: 0, precision: 2 },
				flex: 1,
			},
		],
		[],
	);
	const gridRef = useRef<GridApi<Transaction> | null>(null);

	const handleSave = () => {
		const api = gridRef.current;
		if (!api) return;

		const rows: Transaction[] = [];
		api.forEachNode((n) => n.data && rows.push(n.data));

		save(rows);
	};

	return (
		<div>
			<Text>Successfully parsed transactions</Text>
			<AgGridProvider modules={modules}>
				<div
					style={{
						height: "calc(100vh - 300px)",
					}}
				>
					<AgGridReact
						rowData={transactions}
						defaultColDef={{ flex: 1 }}
						columnDefs={colDefs}
						pagination={true}
						paginationPageSize={20}
						onGridReady={(p) => (gridRef.current = p.api)}
					/>
				</div>
			</AgGridProvider>

			<Group justify="space-between" mt="md">
				<Text size="sm" c="dimmed">
					{transactions.length} transactions · ✓ reconciled
				</Text>
				<Group gap="sm">
					<Button variant="default" onClick={cancel}>
						Cancel
					</Button>
					<Button color="dark" onClick={handleSave}>
						Save
					</Button>
				</Group>
			</Group>
		</div>
	);
}
