import { jsPDF } from 'jspdf';
import type { ProductionReport } from '$lib/types';

export class ReportPDFGenerator {
	private doc: jsPDF;

	constructor() {
		this.doc = new jsPDF();
	}

	private addBackground() {
		// Create a honeycomb pattern background
		this.doc.setFillColor(255, 248, 225); // Light honey color
		this.doc.rect(0, 0, this.doc.internal.pageSize.width, this.doc.internal.pageSize.height, 'F');

		this.doc.setDrawColor(255, 193, 7); // Amber color for honeycomb
		const hexSize = 10;
		for (let x = 0; x < this.doc.internal.pageSize.width + hexSize; x += hexSize * 1.5) {
			for (
				let y = 0;
				y < this.doc.internal.pageSize.height + hexSize;
				y += hexSize * Math.sqrt(3)
			) {
				this.drawHexagon(x + (y % (hexSize * 3) === 0 ? 0 : hexSize * 0.75), y, hexSize);
			}
		}
	}

	private drawHexagon(x: number, y: number, size: number) {
		const angle = Math.PI / 3;
		this.doc.setLineWidth(0.1);
		this.doc.lines(
			[
				[size * Math.cos(angle * 0), size * Math.sin(angle * 0)],
				[size * Math.cos(angle * 1), size * Math.sin(angle * 1)],
				[size * Math.cos(angle * 2), size * Math.sin(angle * 2)],
				[size * Math.cos(angle * 3), size * Math.sin(angle * 3)],
				[size * Math.cos(angle * 4), size * Math.sin(angle * 4)],
				[size * Math.cos(angle * 5), size * Math.sin(angle * 5)]
			],
			x,
			y,
			[1, 1],
			'S'
		);
	}

	private addHeader() {
		this.doc.setFillColor(255, 193, 7);
		this.doc.rect(0, 0, this.doc.internal.pageSize.width, 40, 'F');
		this.doc.setTextColor(0);
		this.doc.setFontSize(24);
		this.doc.text('Apiary Production Report', 20, 25);
		this.addBee(180, 20, 10); // Add a bee to the header
	}

	private addBee(x: number, y: number, size: number) {
		// Draw bee body
		this.doc.setFillColor(255, 193, 7);
		this.doc.ellipse(x, y, size / 2, size / 3, 'F');
		// Draw stripes
		this.doc.setFillColor('#000');
		this.doc.rect(x - size / 4, y - size / 6, size / 6, size / 3, 'F');
		this.doc.rect(x + size / 12, y - size / 6, size / 6, size / 3, 'F');

		// Draw wings
		this.doc.setFillColor(255, 255, 255, 0.7);
		this.doc.ellipse(x - size / 4, y - size / 3, size / 6, size / 8, 'F');
		this.doc.ellipse(x + size / 4, y - size / 3, size / 6, size / 8, 'F');
	}

	private addContent(report: ProductionReport) {
		this.doc.setFontSize(14);
		this.doc.setTextColor(51, 51, 51);
		this.doc.text('Report Details', 20, 50);
		this.doc.line(20, 52, 190, 52);

		this.doc.setFontSize(14);
		const details = [
			{ label: 'Apiary ID:', value: report.apiary_id.toString() },
			{
				label: 'Period:',
				value: `${new Date(report.start_date).toLocaleDateString()} - ${new Date(report.end_date).toLocaleDateString()}`
			},
			{ label: 'Total Honey:', value: `${report.total_honey} kg` },
			{ label: 'Total Expenses:', value: `$${report.total_expenses}` }
		];

		let y = 65;
		details.forEach(({ label, value }) => {
			this.doc.setFont('helvetica', 'bold');
			this.doc.text(label, 20, y);
			this.doc.setFont('helvetica', 'normal');
			this.doc.text(value, 100, y);
			y += 15;
		});

		// Add some flying bees
		for (let i = 0; i < 5; i++) {
			this.addBee(
				Math.random() * this.doc.internal.pageSize.width,
				Math.random() * (this.doc.internal.pageSize.height - 80) + 80,
				5 + Math.random() * 5
			);
		}
	}

	private addFooter() {
		const timestamp = new Date().toLocaleString();
		this.doc.setFontSize(8);
		this.doc.setTextColor(128, 128, 128);
		this.doc.text(`Generated on ${timestamp}`, 20, this.doc.internal.pageSize.height - 10);
		this.doc.text(
			'BeesBiz',
			this.doc.internal.pageSize.width - 40,
			this.doc.internal.pageSize.height - 10
		);
	}

	public generate(report: ProductionReport): void {
		this.addBackground();
		this.addHeader();
		this.addContent(report);
		this.addFooter();

		const formattedDate = new Date().toISOString().split('T')[0];
		this.doc.save(`apiary-${report.apiary_id}-report-${formattedDate}.pdf`);
	}
}
