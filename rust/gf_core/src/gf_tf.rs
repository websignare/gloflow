/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

// use std::io;
// use std::fs;
use std::collections::{HashMap};
use image;
use png;
use png::HasParameters; // needed for png_encoder.set() call

use tensorflow;
use protobuf;
use protobuf::Message;

// these structs are defined in the Rust code files in ./gf_protobuff/ dir,
// which were generated by "protobufc" compiler from the ./gf_protobuff/*.proto definition files.
use crate::gf_protobuff::tf_feature::{Features, Feature, Int64List, BytesList};
use crate::gf_protobuff::tf_example::{Example};

//-------------------------------------------------
// READ
//-------------------------------------------------
#[allow(non_snake_case)]
pub fn read_tf_records(p_target_file_path_str: &str,
    p_img_width_int:  u64,
    p_img_height_int: u64) {


    let input      = std::fs::File::open(p_target_file_path_str).unwrap();
    let input_buff = std::io::BufReader::new(input);

    let mut tf_records_reader = tensorflow::io::RecordReader::new(input_buff);
    let mut tf_example_raw    = [0u8; 3000]; // buffer for individual examples read in from a .tfrecords file

    loop {


        let next = tf_records_reader.read_next(&mut tf_example_raw);


        println!("{:?}", next);


        match next {

            Ok(resp) => match resp {
                Some(len) => { 
                    println!("data received");



                    let data_lst = &tf_example_raw[0..len];
                    

                    // mut - tf_example is mutable here because we're using "taking" (.take_features(), .take_bytes_list(), .take_value()) sub-values
                    //       of a TF Example, and they all require to be operating on a mutable reference of the Example.
                    let mut tf_example = protobuf::parse_from_bytes::<Example>(&data_lst).unwrap();

                    println!("PARSED EXAMPLE");
                    println!("{:?}", tf_example);



                    let tf_features = tf_example.take_features();
                    let mut tf_features_map: ::std::collections::HashMap<::std::string::String, Feature> = tf_features.feature;


                    let tf_feature__label     = tf_features_map.get("label").unwrap();
                    let tf_feature__label_int = tf_feature__label.get_int64_list().get_value()[0];
                    
                    
                    
                    let tf_feature__img = tf_features_map.get_mut("img").unwrap();
                    let tf_feature__img_png_bytes_lst: std::vec::Vec<u8> = tf_feature__img.take_bytes_list().take_value().first().unwrap().to_vec();

                    // println!("{:?}", tf_features_map);
                    // println!("{:?}", tf_feature__label_int);
                    // println!("{:?}", tf_feature__img_png_bytes_lst);


                    // PNG_DECODING

                    // using Cursor because it implements a reader trait which is needed by png::Decoder::new
                    // from Rust docs: "Cursors are typically used with in-memory buffers to allow them to implement Read and/or Write"
                    let tf_feature__img_png_reader_lst: ::std::io::Cursor<std::vec::Vec<u8>> = ::std::io::Cursor::new(tf_feature__img_png_bytes_lst);

                    let png_decoder = png::Decoder::new(tf_feature__img_png_reader_lst);
                    let (info, mut png_reader) = png_decoder.read_info().unwrap();

                    let mut img_buf = vec![0; info.buffer_size()];
                    png_reader.next_frame(&mut img_buf).unwrap();
                    
                    // IMAGE_BUFFER
                    let gf_img_buff: image::ImageBuffer<image::Rgba<u8>, Vec<u8>> = image::ImageBuffer::from_raw(p_img_width_int as u32,
                        p_img_height_int as u32,
                        img_buf).unwrap();

                    
                    println!(" ======= {:?}", gf_img_buff);


                },
                None => break,
            }, 

            Err(tensorflow::io::RecordReadError::CorruptFile) | Err(tensorflow::io::RecordReadError::IoError { .. }) => {
                break;
            }
            _ => {}
        }
    }
}

//-------------------------------------------------
// WRITE
//-------------------------------------------------
// WRITE_TF_RECORDS__FROM_IMG_BUFFER
#[allow(non_snake_case)]
pub fn write_tf_records__from_img_buffer(p_img_buffer: image::ImageBuffer<image::Rgba<u8>, Vec<u8>>,
    p_label_int:      u64,
    p_records_writer: &mut tensorflow::io::RecordWriter<std::io::BufWriter<std::fs::File>>) {


    let mut buf_writer = Vec::new();
    
    // PNG_ENCODING
    {
        let mut png_encoder = png::Encoder::new(&mut buf_writer,
            p_img_buffer.width(),
            p_img_buffer.height());

        png_encoder.set(png::ColorType::RGBA).set(png::BitDepth::Eight);
        
        let mut png_encoder_writer = png_encoder.write_header().unwrap();

        // let data = [255, 0, 0, 255, 0, 0, 0, 255]; // An array containing a RGBA sequence
        let img_data = p_img_buffer.into_raw();
        png_encoder_writer.write_image_data(&img_data).unwrap();
    }

    let img_png_encoded_data_lst: &[u8] = &buf_writer;

    // Vec<Vec<u8>> - used because protobuf::RepeatedField::from_vec() requires a 2D array.
    let img_bytes_lst: Vec<Vec<u8>> = vec![img_png_encoded_data_lst.to_vec()]; // vec![p_img_buffer.into_raw()];
    
    //-----------------
    // FEATURE_LABEL
    let mut tf_feature_label = Feature::new();
    let mut tf_label         = Int64List::new();
    tf_label.set_value(vec![p_label_int as i64]);
    tf_feature_label.set_int64_list(tf_label);

    //-----------------
    // FEATURE_IMG
    let mut tf_feature_img = Feature::new();
    let mut tf_img_bytes   = BytesList::new();

    tf_img_bytes.set_value(protobuf::RepeatedField::from_vec(img_bytes_lst));
    tf_feature_img.set_bytes_list(tf_img_bytes);

    //-----------------
    // FEATURES
    
    let mut tf_feature_map = HashMap::new();
    tf_feature_map.insert("label".to_string(), tf_feature_label);
    tf_feature_map.insert("img".to_string(),   tf_feature_img);

    let mut tf_features = Features::new();
    tf_features.set_feature(tf_feature_map);

    //-----------------
    // EXAMPLE
    let mut tf_example = Example::new();
    tf_example.set_features(tf_features);

    let tf_example_bytes_lst = tf_example.write_to_bytes().unwrap();

    //-----------------
    // TF_RECORDS_WRITER
    p_records_writer.write_record(&tf_example_bytes_lst).unwrap();

    //-----------------
}

//-------------------------------------------------
// WRITE_TF_RECORDS__TO_FILE
pub fn write_tf_records__to_file(p_output_file_path_str: &str) {
    
    let label_int         = 0 as u64;
    let img_file_path_str = "data/output_ml/generated/train/rect/test-rect-0.png";

    // TF_RECORDS_WRITER
    let mut tf_records_writer = get_tf_records__writer(p_output_file_path_str);

    //-----------------
    // IMAGE_BUFFER
    
    let img:         image::DynamicImage                          = image::open(img_file_path_str).unwrap();
    let gf_img_buff: image::ImageBuffer<image::Rgba<u8>, Vec<u8>> = img.to_rgba();

    // WRITE_TF_RECORD
    write_tf_records__from_img_buffer(gf_img_buff,
        label_int,
        &mut tf_records_writer);

    for x in 0..10 {

        //-----------------
        // WRITE_RECORD
        
        println!("===");
        // record_writer.write_record(&tf_example_bytes_lst).unwrap();

        //-----------------
    }
}

//-------------------------------------------------
// GET_TF_RECORDS__WRITER
pub fn get_tf_records__writer(p_output_file_path_str: &str) -> tensorflow::io::RecordWriter<std::io::BufWriter<std::fs::File>> {

    let f = ::std::fs::OpenOptions::new()
        .write(true)
        // IMPORTANT!! - truncate() - this is critical for the generated .tfrecords file not to be corrupted
        //                            (if its overwriting an existing file).
        //                            if this is set to false, and an .tfrecords file with the same name already exists,
        //                            this will overwrite the file but if the data is smaller in size then the existing file 
        //                            it will leave the old data as padding... this will result in a corrupted file and 
        //                            TensorFlow will throw an exception when parsing it.
        .truncate(true)
        .create(true)
        .open(p_output_file_path_str)
        .unwrap();
    
    // buffered writer for the final output .tfrecords file
    let tf_record_writer = tensorflow::io::RecordWriter::new(::std::io::BufWriter::new(f));
    return tf_record_writer;
}