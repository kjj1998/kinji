interface deltaProps {
	current: number;
	previous: number;
}

export function calculateDelta({ current, previous }: deltaProps): number | null {
	if (previous === 0) return null;
	return ((current - previous) / previous) * 100;
}
