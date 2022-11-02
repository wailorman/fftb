class FindFreeTask < ApplicationInteractor
  object :performer, class: Performer

  def execute
    Task.transaction do
      Task.with_advisory_lock('find_free_task', transaction: true) do
        found = Task.published.not_occupied.first

        return nil unless found

        found.occupied_at = Time.current
        found.occupied_by = performer

        unless found.save
          errors.merge!(found.errors)
          return nil
        end

        found
      end
    end
  end
end
