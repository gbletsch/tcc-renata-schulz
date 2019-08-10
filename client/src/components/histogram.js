import React, { PureComponent } from 'react'
import {
  BarChart, Bar, Cell, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer
} from 'recharts'

export default class Histogram extends PureComponent {
  constructor (props) {
    super(props)
    this.state = {
      y_scale: 'linear'
    }

    this.onChange = event => {
      const state = Object.assign({}, this.state)
      const field = event.target.name
      state[field] = event.target.value
      this.setState(state)
    }
  }

  render () {

    var yScale = this.props.yScale

    return (

      <div style={{ width: '100%', height: 300 }}>

        <ResponsiveContainer>

          <BarChart
            data={this.props.data[0]}
            margin={{
              top: 5, right: 30, left: 20, bottom: 5
            }}
          >
            <CartesianGrid strokeDasharray='3 3' />
            <XAxis dataKey='bin' />
            <YAxis scale={yScale} domain={[0.01, 'auto']} allowDataOverflow />
            <Tooltip />
            <Legend />
            <Bar dataKey='count' fill='#8884d8' name='Número de internações' />
          </BarChart>
        </ResponsiveContainer>
      </div>
    )
  }
}
