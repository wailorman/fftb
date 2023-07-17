class FileSelections::ItemsController < ApplicationController
  before_action :set_file_selection
  before_action :set_file_selection_item

  def reveal
    @file_selection_item.update!(removed: false)
    redirect_to @file_selection
  end

  def hide
    @file_selection_item.update!(removed: true)
    redirect_to @file_selection
  end

  private def set_file_selection
    @file_selection = FileSelection.find(params[:file_selection_id])
  end

  private def set_file_selection_item
    @file_selection_item = @file_selection.items.unscoped.find(params[:id])
  end
end
