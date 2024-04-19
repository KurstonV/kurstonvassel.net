package postgresql

import (
	"database/sql"
	"errors"

	"kurstonvassel.net/quotebox/pkg/models"
)

type EmployeeModel struct {
	DB *sql.DB
}

func (m *EmployeeModel) Insert(employee, email, job, salary string) (int, error) {
	var id int

	s := `
	    INSERT INTO employees(employee_name, email, job, salary)
	    VALUES ($1, $2, $3, $4)
		RETURNING employee_id
	`
	err := m.DB.QueryRow(s, employee, email, job, salary).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *EmployeeModel) Read() ([]*models.Employee, error) {
	s := `
		SELECT employee_name, email, job, salary
		FROM employees
		LIMIT 10
	`
	rows, err := m.DB.Query(s)
	if err != nil {
		return nil, err
	}
	// cleanup before we leave Read()
	defer rows.Close()

	quotes := []*models.Employee{}

	for rows.Next() {
		q := &models.Employee{}
		err = rows.Scan(&q.Author_name, &q.Category, &q.Body)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, q)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return quotes, nil
}

func (m *EmployeeModel) Get(id int) (*models.Employee, error) {
	s := `
		SELECT employee_name, email, job, salary
		FROM employees
		WHERE employees_id = $1
	`
	q := &models.Employee{}
	err := m.DB.QueryRow(s, id).Scan(&e.Employee_name, &e.Email, &e.Job, &e.Salary)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		} else {
			return nil, err
		}
	}
	return q, nil
}
