class OrdersController < ApplicationController
  PUBLISH_COMMIT = 'Publish'.freeze
  REBUILD_COMMIT = 'Rebuild tasks'.freeze
  ALLOWED_ORDER_TYPES = [
    Orders::Convert
  ].map(&:name).freeze

  before_action :set_order, except: %i[index choose_type new create]

  def index
    @orders = Order.page(params[:page]).per(30)
  end

  def choose_type
  end

  def new
    @order = Order.new(permitted_params)

    if @order.type.blank?
      return redirect_to choose_type_orders_path
    end

    @order.file_selection = current_file_selection
    @order.set_default_values
  end

  def create
    @order = order_class.new(permitted_params)
    @order.file_selection = current_file_selection

    result = apply_changes(build_tasks: true) do
      @order.set_default_values
    end

    if result
      redirect_to edit_order_path(@order)
    else
      render :new
    end
  end

  def show
  end

  def edit
    if @order.state.published?
      redirect_to order_path(@order)
      return
    end

    render :edit
  end

  def update
    result = apply_changes(build_tasks: params[:commit] == REBUILD_COMMIT) do
      @order.publish if params[:commit] == PUBLISH_COMMIT
    end

    if result
      redirect_to edit_order_path(@order)
    else
      render :edit
    end
  end

  def cancel
    @order.cancel

    unless @order.save
      flash[:alert] = "Failed to cancel order: #{human_errors(@order.errors)}"
    end

    redirect_to action: :show
  end

  private def set_order
    @order = Order.find(params[:id])
  end

  # rubocop:disable Style/SymbolArray
  private def permitted_params
    f_params = params.fetch(:order, {}).permit(:type)
    is_convert_order = params.fetch(:orders_convert, {}).permit!.present?

    if is_convert_order
      params
        .require(:orders_convert)
        .permit(
          payload_attributes: [
            :id,
            :video_muxer,
            :video_opts,
            :audio_muxer,
            :audio_opts,
            :output_rclone_path
          ],
          tasks_attributes: [
            :id,
            :_destroy,
            payload_attributes: [
              :id,
              :input_rclone_path,
              :output_rclone_path,
              :string_opts
            ]
          ]
        ).tap do |p_params|
          p_params[:type] = Orders::Convert.name
          p_params[:file_selection_id] = current_file_selection&.id
        end
    else
      f_params
    end
  end
  # rubocop:enable Style/SymbolArray

  private def order_class
    (([permitted_params[:type]] & ALLOWED_ORDER_TYPES).first || 'Order').constantize
  end

  private def apply_changes(build_tasks: false)
    result = true

    Order.transaction do
      @order.assign_attributes(permitted_params)

      yield if block_given?

      unless @order.save
        raise ActiveRecord::Rollback
      end

      if build_tasks
        build_service = BuildConvertOrderService.new(@order)

        unless build_service.perform
          @order.errors.add(:base, "Failed to build tasks: #{build_service.errors.full_messages.join(', ')}")
          raise ActiveRecord::Rollback
        end
      end
    rescue ActiveRecord::Rollback => e
      result = false
      raise e
    end

    result
  end
end
