import React, { Component } from 'react'
import Downshift from 'downshift'

export default class downshift extends Component {
  render () {
    const items = this.props.options

    return (
      <Downshift
        onChange={selection => alert(
          selection ? `You selected ${selection.value}` : 'Selection Cleared'
        )}
        itemToString={item => (item ? item.value : '')}
      >
        {({
          getInputProps,
          getItemProps,
          getLabelProps,
          getMenuProps,
          isOpen,
          inputValue,
          highlightedIndex,
          selectedItem
        }) => (
          <div>
            <label {...getLabelProps()}>Enter a fruit</label>
            <input {...getInputProps()} />
            <ul {...getMenuProps()}>
              {isOpen
                ? items
                  .filter(item => !inputValue || item.value.includes(inputValue))
                  .map((item, index) => (
                    <li
                      {...getItemProps({
                        key: item.value,
                        index,
                        item,
                        style: {
                          backgroundColor:
                          highlightedIndex === index ? 'lightgray' : 'white',
                          fontWeight: selectedItem === item ? 'bold' : 'normal'
                        }
                      })}
                    >
                      {item.value}
                    </li>
                  ))
                : null}
            </ul>
          </div>
        )}
      </Downshift>
    )
  }
}
