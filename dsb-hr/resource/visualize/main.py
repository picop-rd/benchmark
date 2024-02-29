import pandas as pd
import matplotlib.pyplot as plt
import argparse

def create_stacked_bar_chart(data_dir):
	for i in ("cpu", "memory"):
		plt.clf()
		# Load the data from the CSV file
		df = pd.read_csv(data_dir+"/"+i+".csv", index_col='type')
		print(df)
		print(df.index)
		print(df.columns)
		print(df.loc['base'])
		print(df.loc['picop'])

		plt.bar(df.columns, df.loc['base'], label='base', alpha=0.5)
		plt.bar(df.columns, df.loc['picop'], label='picop', alpha=0.5)
		plt.legend(loc='lower left')
		plt.grid(axis='y')

		# Adding labels and title
		plt.xlabel('Senario')
		if i == "cpu":
			plt.ylabel('vCPU [m core]')
		else:
			plt.ylabel('Memory [MiB]')

		# Save the plot to the specified output path
		plt.savefig(data_dir+"/"+i+".png")

# Example usage
# create_stacked_bar_chart('path_to_input.csv', 'path_to_output.png')

def main():
	parser = argparse.ArgumentParser()
	parser.add_argument('data_dir', type=str)
	args = parser.parse_args()
	create_stacked_bar_chart(args.data_dir)

if __name__ == '__main__':
	main()
