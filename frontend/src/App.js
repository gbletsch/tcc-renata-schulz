import React, { Component } from 'react'
import Histogram from './components/histogram'
import LinePlot from './components/lineplot'

export default class App extends Component {
  constructor (props) {
    super(props)
    this.state = {
      items: [],
      isLoaded: false,
      cid: '',
      uf: '',
      ano: ''
    }

    this.onSave = event => {
      event.preventDefault()
      const path = `http://localhost:8080/index?CID=${this.state.cid}&UF=${this.state.uf}&Ano=${this.state.ano}`

      fetch(path)
        .then(res => res.json())
        .then(json => {
          this.setState({
            isLoaded: true,
            items: json
          })
          this.setState({
            isLoaded: true,
            items: json
          })
        })
    }

    this.onChange = event => {
      const state = Object.assign({}, this.state)
      const field = event.target.name
      state[field] = event.target.value
      this.setState(state)
    }
  }

  componentDidMount () {
    const path = `http://localhost:8080/index?CID=${this.state.cid}&UF=${this.state.uf}&Ano=${this.state.ano}`
    fetch(path)
      .then(res => res.json())
      .then(json => {
        this.setState({
          isLoaded: true,
          items: json
        })
      })
  }

  render () {
    var isLoaded = this.state.isLoaded
    var items = this.state.items
    var options = items.map(opt => opt['options'])
    var cids = options.map(cid => cid['cids'])
    var ufs = options.map(uf => uf['ufs'])
    var years = options.map(y => y['years'])
    var uss = items.map(u => u['USS_hist'])
    var idade = items.map(i => i['age_hist'])
    var lineplot = items.map(lp => lp['lineplots'])

    if (!isLoaded) {
      return <div>Loading...</div>
    } else {
      return (
        <div className='container'>
          <form className='form-group' onSubmit={this.onSave}>
            <div class='row'>
              <div class='col-sm'>
                <select className='form-control' name='cid' value={this.state.cid} onChange={this.onChange}>
                  <option value='Todos'>CID</option>
                  {cids[0].map(item =>
                    <option value={item} >{item}</option>
                  )}
                </select>
              </div>
              <div class='col-sm'>

                <select className='form-control' name='uf' value={this.state.uf} onChange={this.onChange}>
                  <option value='Todos'>UF</option>
                  {ufs[0].map(item =>
                    <option value={item} >{item}</option>
                  )}
                </select>
              </div>
              <div class='col-sm'>
                <select className='form-control' name='ano' value={this.state.ano} onChange={this.onChange}>
                  <option value='Todos'>Ano</option>
                  {years[0].map(item =>
                    <option value={item} >{item}</option>
                  )}
                </select>
              </div>
              <div class='col-sm'>
                <button type='submit' onClick={this.onSave} >Submit</button>
              </div>
            </div>
          </form>
          <div>
            <h2>Evolução das internações e mortalidade</h2>
            <LinePlot data={lineplot} />
          </div>

          <div>
            <h2>Histograma Idade</h2>
            <Histogram data={idade} yScale='linear' />
          </div>

          <div>
            <h2>Histograma Repasse Total (US$)</h2>
            <p>Escala logaritmica</p>
            <Histogram data={uss} yScale='log' />
          </div>
        </div>
      )
    }
  }
}
