import React, { PureComponent } from 'react'
import {
  LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer
} from 'recharts'

export default class LinePlot extends PureComponent {
  static jsfiddleUrl = 'https://jsfiddle.net/alidingling/zjb47e83/';

  render () {
    return (
      <div style={{ width: '100%', height: 300 }}>

        <ResponsiveContainer>
          <LineChart
            data={this.props.data[0]}
            margin={{
              top: 5, right: 30, left: 20, bottom: 5
            }}
          >
            <CartesianGrid strokeDasharray='3 3' />
            <XAxis dataKey='year' />
            <YAxis yAxisId='left' />
            <YAxis yAxisId='right' orientation='right' />
            <Tooltip />
            <Legend />
            <Line yAxisId='left' type='monotone' dataKey='admissions' name='Internações' stroke='#8884d8' activeDot={{ r: 8 }} />
            <Line yAxisId='right' type='monotone' dataKey='mortality' name='Mortalidade' stroke='#82ca9d' />
          </LineChart>
        </ResponsiveContainer>
        </div>
    )
  }
}
