package testing

import (
	context "context"
	driver "database/sql/driver"
)

type MockDriver struct {
	Conn driver.Conn
}

func (m *MockDriver) Open(name string) (driver.Conn, error) {
	if m.Conn == nil {
		return nil, driver.ErrBadConn
	}
	return m.Conn, nil
}

type MockConn struct {
	Stmt driver.Stmt
	Err  error
}

func (m *MockConn) Begin() (driver.Tx, error) {
	return nil, driver.ErrBadConn
}

func (m *MockConn) Close() error {
	return nil
}

func (m *MockConn) Prepare(query string) (driver.Stmt, error) {
	return m.Stmt, m.Err
}

func (m *MockConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return m.Stmt, m.Err
}

type MockStmt struct {
	count int
	Err   []error
	Rows  []*StubRows
}

func (m *MockStmt) Close() error {
	return nil
}

func (m *MockStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, driver.ErrBadConn
}

func (m *MockStmt) NumInput() int {
	return 0
}

func (m *MockStmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, nil
}

func (m *MockStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if m.count >= len(m.Rows) {
		m.count = 0
	}
	rows := m.Rows[m.count]
	err := m.Err[m.count]
	m.count++

	return rows, err
}

type StubRows struct{}

func (m *StubRows) Columns() []string {
	return []string{}
}
func (m *StubRows) Close() error {
	return nil
}
func (m *StubRows) Next(dest []driver.Value) error {
	return nil
}
