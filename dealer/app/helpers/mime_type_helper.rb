module MimeTypeHelper
  def mime_type_icon(mime_type)
    case mime_type
    when 'inode/directory'
      'bi bi-folder2'
    when %r{^video/}
      'bi bi-film'
    when %r{^audio/}
      'bi bi-music-note-beamed'
    end
  end

  def media_mime_type?(mime_type)
    video_mime_type?(mime_type) || audio_mime_type?(mime_type)
  end

  def video_mime_type?(mime_type)
    mime_type.start_with?('video/')
  end

  def audio_mime_type?(mime_type)
    mime_type.start_with?('audio/')
  end
end
