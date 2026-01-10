// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package models

import (
	"testing"

	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRelation_Create(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindSubtask,
		}
		err := rel.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		err = s.Commit()
		require.NoError(t, err)
		db.AssertExists(t, "task_relations", map[string]interface{}{
			"task_id":       1,
			"other_task_id": 2,
			"relation_kind": RelationKindSubtask,
			"created_by_id": 1,
		}, false)
	})
	t.Run("Two Tasks In Different Projects", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       1,
			OtherTaskID:  13,
			RelationKind: RelationKindSubtask,
		}
		err := rel.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		err = s.Commit()
		require.NoError(t, err)
		db.AssertExists(t, "task_relations", map[string]interface{}{
			"task_id":       1,
			"other_task_id": 13,
			"relation_kind": RelationKindSubtask,
			"created_by_id": 1,
		}, false)
	})
	t.Run("Already Existing", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       1,
			OtherTaskID:  29,
			RelationKind: RelationKindSubtask,
		}
		err := rel.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrRelationAlreadyExists(err))
	})
	t.Run("Same Task", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:      1,
			OtherTaskID: 1,
		}
		err := rel.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrRelationTasksCannotBeTheSame(err))
	})
	t.Run("cycle with one subtask", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       29,
			OtherTaskID:  1,
			RelationKind: RelationKindSubtask,
		}
		err := rel.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with multiple subtasks", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindSubtask,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  3,
			RelationKind: RelationKindSubtask,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel3 := TaskRelation{
			TaskID:       3,
			OtherTaskID:  4,
			RelationKind: RelationKindSubtask,
		}
		err = rel3.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Cycle happens here
		rel4 := TaskRelation{
			TaskID:       4,
			OtherTaskID:  2,
			RelationKind: RelationKindSubtask,
		}
		err = rel4.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with multiple subtasks tasks and relation back to parent", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindSubtask,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  3,
			RelationKind: RelationKindSubtask,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel3 := TaskRelation{
			TaskID:       3,
			OtherTaskID:  4,
			RelationKind: RelationKindSubtask,
		}
		err = rel3.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Cycle happens here
		rel4 := TaskRelation{
			TaskID:       4,
			OtherTaskID:  1,
			RelationKind: RelationKindSubtask,
		}
		err = rel4.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with one parenttask", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       1,
			OtherTaskID:  29,
			RelationKind: RelationKindParenttask,
		}
		err := rel.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with multiple parenttasks", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindParenttask,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  3,
			RelationKind: RelationKindParenttask,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel3 := TaskRelation{
			TaskID:       3,
			OtherTaskID:  4,
			RelationKind: RelationKindParenttask,
		}
		err = rel3.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Cycle happens here
		rel4 := TaskRelation{
			TaskID:       4,
			OtherTaskID:  2,
			RelationKind: RelationKindParenttask,
		}
		err = rel4.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with multiple parenttasks and relation back to parent", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindParenttask,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  3,
			RelationKind: RelationKindParenttask,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel3 := TaskRelation{
			TaskID:       3,
			OtherTaskID:  4,
			RelationKind: RelationKindParenttask,
		}
		err = rel3.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Cycle happens here
		rel4 := TaskRelation{
			TaskID:       4,
			OtherTaskID:  1,
			RelationKind: RelationKindParenttask,
		}
		err = rel4.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with one follows relation", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		// Create: Task 1 follows Task 2
		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindFollows,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Try to create: Task 2 follows Task 1 (would create a cycle)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  1,
			RelationKind: RelationKindFollows,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with multiple follows relations", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		// Create chain: 1 follows 2 follows 3 follows 4
		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindFollows,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  3,
			RelationKind: RelationKindFollows,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel3 := TaskRelation{
			TaskID:       3,
			OtherTaskID:  4,
			RelationKind: RelationKindFollows,
		}
		err = rel3.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Try to create: 4 follows 2 (would create a cycle)
		rel4 := TaskRelation{
			TaskID:       4,
			OtherTaskID:  2,
			RelationKind: RelationKindFollows,
		}
		err = rel4.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with one precedes relation", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		// Create: Task 1 precedes Task 2
		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindPreceeds,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Try to create: Task 2 precedes Task 1 (would create a cycle)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  1,
			RelationKind: RelationKindPreceeds,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with multiple precedes relations", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		// Create chain: 1 precedes 2 precedes 3 precedes 4
		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindPreceeds,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  3,
			RelationKind: RelationKindPreceeds,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel3 := TaskRelation{
			TaskID:       3,
			OtherTaskID:  4,
			RelationKind: RelationKindPreceeds,
		}
		err = rel3.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Try to create: 4 precedes 1 (would create a cycle back to the start)
		rel4 := TaskRelation{
			TaskID:       4,
			OtherTaskID:  1,
			RelationKind: RelationKindPreceeds,
		}
		err = rel4.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with one blocking relation", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		// Create: Task 1 is blocking Task 2
		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindBlocking,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Try to create: Task 2 is blocking Task 1 (would create a cycle)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  1,
			RelationKind: RelationKindBlocking,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with multiple blocking relations", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		// Create chain: 1 blocking 2 blocking 3 blocking 4
		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindBlocking,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  3,
			RelationKind: RelationKindBlocking,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel3 := TaskRelation{
			TaskID:       3,
			OtherTaskID:  4,
			RelationKind: RelationKindBlocking,
		}
		err = rel3.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Try to create: 4 blocking 2 (would create a cycle)
		rel4 := TaskRelation{
			TaskID:       4,
			OtherTaskID:  2,
			RelationKind: RelationKindBlocking,
		}
		err = rel4.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with one blocked relation", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		// Create: Task 1 is blocked by Task 2
		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindBlocked,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Try to create: Task 2 is blocked by Task 1 (would create a cycle)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  1,
			RelationKind: RelationKindBlocked,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
	t.Run("cycle with multiple blocked relations", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		// Create chain: 1 blocked by 2 blocked by 3 blocked by 4
		rel1 := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindBlocked,
		}
		err := rel1.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel2 := TaskRelation{
			TaskID:       2,
			OtherTaskID:  3,
			RelationKind: RelationKindBlocked,
		}
		err = rel2.Create(s, &user.User{ID: 1})
		require.NoError(t, err)
		rel3 := TaskRelation{
			TaskID:       3,
			OtherTaskID:  4,
			RelationKind: RelationKindBlocked,
		}
		err = rel3.Create(s, &user.User{ID: 1})
		require.NoError(t, err)

		// Try to create: 4 blocked by 1 (would create a cycle back to the start)
		rel4 := TaskRelation{
			TaskID:       4,
			OtherTaskID:  1,
			RelationKind: RelationKindBlocked,
		}
		err = rel4.Create(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskRelationCycle(err))
	})
}

func TestTaskRelation_Delete(t *testing.T) {
	u := &user.User{ID: 1}

	t.Run("Normal", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       1,
			OtherTaskID:  29,
			RelationKind: RelationKindSubtask,
		}
		err := rel.Delete(s, u)
		require.NoError(t, err)
		err = s.Commit()
		require.NoError(t, err)
		db.AssertMissing(t, "task_relations", map[string]interface{}{
			"task_id":       1,
			"other_task_id": 29,
			"relation_kind": RelationKindSubtask,
		})
		db.AssertMissing(t, "task_relations", map[string]interface{}{
			"task_id":       29,
			"other_task_id": 1,
			"relation_kind": RelationKindParenttask,
		})
	})
	t.Run("Not existing", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       9999,
			OtherTaskID:  3,
			RelationKind: RelationKindSubtask,
		}
		err := rel.Delete(s, u)
		require.Error(t, err)
		assert.True(t, IsErrRelationDoesNotExist(err))
	})
}

func TestTaskRelation_CanCreate(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       1,
			OtherTaskID:  2,
			RelationKind: RelationKindSubtask,
		}
		can, err := rel.CanCreate(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.True(t, can)
	})
	t.Run("Two tasks on different projects", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       1,
			OtherTaskID:  32,
			RelationKind: RelationKindSubtask,
		}
		can, err := rel.CanCreate(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.True(t, can)
	})
	t.Run("No update permissions on base task", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       14,
			OtherTaskID:  1,
			RelationKind: RelationKindSubtask,
		}
		can, err := rel.CanCreate(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.False(t, can)
	})
	t.Run("No update permissions on base task, but read permissions", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       15,
			OtherTaskID:  1,
			RelationKind: RelationKindSubtask,
		}
		can, err := rel.CanCreate(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.False(t, can)
	})
	t.Run("No read permissions on other task", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       1,
			OtherTaskID:  14,
			RelationKind: RelationKindSubtask,
		}
		can, err := rel.CanCreate(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.False(t, can)
	})
	t.Run("Nonexisting base task", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       999999,
			OtherTaskID:  1,
			RelationKind: RelationKindSubtask,
		}
		can, err := rel.CanCreate(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskDoesNotExist(err))
		assert.False(t, can)
	})
	t.Run("Nonexisting other task", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		rel := TaskRelation{
			TaskID:       1,
			OtherTaskID:  999999,
			RelationKind: RelationKindSubtask,
		}
		can, err := rel.CanCreate(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrTaskDoesNotExist(err))
		assert.False(t, can)
	})
}
