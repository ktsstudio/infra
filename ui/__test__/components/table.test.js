import { render, screen, fireEvent } from '@testing-library/react'
import Table from '../../components/table'

const columns = [
  { accessorKey: 'name', header: 'Name', cell: info => info.getValue() },
]

describe('Table', () => {
  it('renders API rows and requests the next server page', () => {
    const onPageChange = vi.fn()
    render(
      <Table
        columns={columns}
        data={[{ id: 'one', name: 'Cluster one' }]}
        count={2}
        pageSize={1}
        pageCount={2}
        onPageChange={onPageChange}
      />
    )
    expect(screen.getByText('Cluster one')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Prev' })).toBeDisabled()
    fireEvent.click(screen.getByRole('button', { name: 'Next' }))
    expect(onPageChange).toHaveBeenCalledWith({ pageIndex: 1, pageSize: 1 })
  })

  it('keeps checkbox selection using original API IDs', () => {
    const setSelectedRowIds = vi.fn()
    render(
      <Table
        columns={columns}
        data={[{ id: 'api-id', name: 'Cluster' }]}
        allowDelete
        setSelectedRowIds={setSelectedRowIds}
      />
    )
    fireEvent.click(screen.getAllByRole('checkbox')[1])
    expect(setSelectedRowIds).toHaveBeenLastCalledWith(['api-id'])
  })
})
