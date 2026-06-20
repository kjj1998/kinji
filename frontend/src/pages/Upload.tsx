import { Box, Stack, Text } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useRef, useState } from "react";
import {
	ConfirmStatement,
	Header,
	ParsingStatement,
	ReviewStatement,
	SelectStatement,
} from "../components";
import { useSaveTransactions } from "../hooks";
import { importStatement } from "../services";
import type { Transaction } from "../types";

interface UploadProps {
	userId: string;
}

const MAX_SIZE = 10 * 1024 ** 2;

type Phase = "select" | "parsing" | "review" | "error";
export type Stage = "uploaded" | "validating" | "parsing" | "checking_balances";

export function Upload({ userId }: UploadProps) {
	const [stage, setStage] = useState<Stage | null>(null);
	const [phase, setPhase] = useState<Phase>("select");
	const [file, setFile] = useState<File | null>(null);
	const [password, setPassword] = useState("");
	const [transactions, setTransactions] = useState<Transaction[]>([]);
	const [_errorMsg, setErrorMsg] = useState("");

	const abortRef = useRef<AbortController | null>(null);

	const save = useSaveTransactions(userId);

	const handleSave = (rows: Transaction[]) => {
		save.mutate(rows, {
			onSuccess: (data) => {
				notifications.show({
					message: `✓ ${data.length} transactions saved`,
					withCloseButton: true,
				});
				clearFile();
				setTransactions([]);
			},
			onError: (err) =>
				notifications.show({
					color: "red",
					message: err.message,
					withCloseButton: true,
				}),
		});
	};

	const handleConfirm = async () => {
		if (!file) return;
		setStage(null);
		setErrorMsg("");
		setPhase("parsing");

		const controller = new AbortController();
		abortRef.current = controller;

		try {
			await importStatement(
				file,
				userId,
				password,
				{
					onProgress: setStage,
					onDone: (txns) => {
						setTransactions(txns);
						setPhase("review");
					},
					onError: (msg) => {
						setErrorMsg(msg);
						setPhase("error");
					},
				},
				controller.signal,
			);
		} catch (err) {
			if (err instanceof DOMException && err.name === "AbortError") return; // cancelled
			setErrorMsg("Something went wrong while uploading.");
			setPhase("error");
		} finally {
			abortRef.current = null;
		}
	};

	const clearFile = () => {
		abortRef.current?.abort();
		abortRef.current = null;
		setPhase("select");
		setFile(null);
		setPassword("");
		setStage(null);
		setErrorMsg("");
	};

	return (
		<>
			<Header text="Upload statement" />
			<Box px="md" pt="lg">
				<Stack gap="lg">
					<Text size="sm" c="dimmed">
						Drop a bank or card statement — we'll read it for you
					</Text>
					{phase === "parsing" ? (
						<ParsingStatement file={file} stage={stage} onCancel={clearFile} />
					) : phase === "review" ? (
						<ReviewStatement
							transactions={transactions}
							save={handleSave}
							cancel={clearFile}
						/>
					) : phase === "error" ? (
						<div>Error</div>
					) : file ? (
						<ConfirmStatement
							file={file}
							password={password}
							onPasswordChange={setPassword}
							onCancel={clearFile}
							onConfirm={handleConfirm}
						/>
					) : (
						<SelectStatement maxSize={MAX_SIZE} onFileSelect={setFile} />
					)}
				</Stack>
			</Box>
		</>
	);
}
